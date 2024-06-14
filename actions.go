package main

import (
	"context"
	"encoding/json"
	"io"
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
		{"email", usuarioLogin.Email},
		{"password", usuarioLogin.Password},
	}

	resultadoBusca, erroNaBusca := database.FindOneUser(filter)

	if erroNaBusca != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	token, err := utility.GenerateToken(resultadoBusca.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome, " + resultadoBusca.Username)
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

	results, err := database.FindAllUsers(client, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		// log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		// log.Fatal(err)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	// defer file.Close()

	fileDescription := r.FormValue("description")

	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, handler.Filename)

	// Cria o arquivo temporário
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// defer tempFile.Close()

	client := awsfunctions.Set()

	result, errUpload := awsfunctions.UploadObject(client, "projeto-ltp2", handler.Filename, tempFile)
	if errUpload != nil {
		http.Error(w, errUpload.Error(), http.StatusInternalServerError)
		return
	}

	image := models.Imagem{
		Id:          primitive.NewObjectID(),
		Filename:    handler.Filename,
		FileUrl:     result,
		Description: fileDescription,
	}

	resultado, err := database.InsertOneImage(image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado.InsertedID)
}
