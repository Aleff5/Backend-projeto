package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	awsfunctions "projetov2/Backend-projeto/AwsFunctions"
	database "projetov2/Backend-projeto/Database"
	"projetov2/Backend-projeto/models"
	"projetov2/Backend-projeto/utility"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var usuarioLogin models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuarioLogin)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filter := bson.D{
		{Key: "email", Value: usuarioLogin.Email},
		{Key: "password", Value: usuarioLogin.Password},
	}

	resultadoBusca, erroNaBusca := database.FindOneUser(filter)

	if erroNaBusca != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	token, err := utility.GenerateToken(resultadoBusca.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome, " + resultadoBusca.Username)

	log.Println("Cookie 'token' set with value:", token)
}

func Singup(w http.ResponseWriter, r *http.Request) {
	var novoUsuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&novoUsuario)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := database.InsertOneUser(novoUsuario)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func AdminView(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectBd()
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Buscar todos os usuários retorna lista de primitive.D
	users, err := database.FindAllUsers(client, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Buscar URLs das imagens retorna lista de strings
	imageUrls, err := database.FindUrl() // Suponha que esta função retorna os URLs das imagens
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Criar estrutura para enviar como resposta
	responseData := map[string]interface{}{
		"users":      users,
		"image_urls": imageUrls,
	}

	// Converter para JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Definir o cabeçalho Content-Type e enviar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Define o cookie com uma data de expiração passada para removê-lo
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Altere para true se estiver usando HTTPS
		SameSite: http.SameSiteDefaultMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println("Error getting file from form:", err)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	defer file.Close()

	fileName := r.FormValue("fileName")
	description := r.FormValue("description")

	log.Println("File name:", fileName)            // Log do nome do arquivo
	log.Println("Description:", description)       // Log da descrição
	log.Println("File handler:", handler.Filename) // Log do nome original do arquivo

	if fileName == "" || description == "" || handler.Filename == "" {
		log.Println("Missing form fields")
		http.Error(w, "Missing form fields", http.StatusBadRequest)
		return
	}

	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fileName+filepath.Ext(handler.Filename))

	// Cria o arquivo temporário
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Println("Error creating temp file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Println("Error copying file to temp file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reabrir o arquivo temporário para fazer upload para o S3
	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		log.Println("Error reopening temp file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	client := awsfunctions.Set()

	result, errUpload := awsfunctions.UploadObject(client, "projeto-ltp2", fileName+filepath.Ext(handler.Filename), tempFile)
	if errUpload != nil {
		log.Println("Error uploading to S3:", errUpload)
		http.Error(w, errUpload.Error(), http.StatusInternalServerError)
		return
	}

	image := models.Imagem{
		Id:          primitive.NewObjectID(),
		Filename:    fileName + filepath.Ext(handler.Filename),
		FileUrl:     result,
		Description: description,
	}

	resultado, err := database.InsertOneImage(image)
	if err != nil {
		log.Println("Error inserting image to database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado.InsertedID)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	var imagem models.Imagem
	err := json.NewDecoder(r.Body).Decode(&imagem)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	filename := imagem.Filename

	client := awsfunctions.Set()
	erroDelete := awsfunctions.DeleteObject(client, "projeto-ltp2", filename)
	if erroDelete != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	filter := bson.D{
		{Key: "filename", Value: filename},
	}
	errDb := database.DeleteImage(filter)
	if errDb != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(imagem.Filename + " successfully deleted")
}

// func ShowAll(w http.ResponseWriter, r *http.Request) {
// 	imagens, err := database.FindAllImages()
// 	if err != nil {
// 		http.Error(w, "Internal Error", http.StatusInternalServerError)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(imagens)
// }

func ImageGen(w http.ResponseWriter, r *http.Request) {
	nome, descricao, url, _ := utility.GenerateRandomImage()
	dados := []string{nome, descricao, url}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dados)
}

// função de enviar todas as imagens (admin)
func teste(w http.ResponseWriter, r *http.Request) {
	resultado, _ := database.FindUrl()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}
