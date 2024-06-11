package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"projetov2/Backend-projeto/models"
	"projetov2/Backend-projeto/utility"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message)
}

func Login(w http.ResponseWriter, r *http.Request) {
	client := ConnectBd()
	defer client.Disconnect(context.Background())

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

	collection := client.Database("ProjetoLTP2").Collection("Usuarios")
	resultadoBusca := models.Usuario{}

	erroNaBusca := collection.FindOne(context.Background(), filter).Decode(&resultadoBusca)
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

	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	result, err := collection.InsertOne(context.Background(), novoUsuario)
	if err != nil {
		log.Printf("Failed to create user: %v\n", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func AdminView(w http.ResponseWriter, r *http.Request) {
	client := ConnectBd()
	defer client.Disconnect(context.Background())

	collection := client.Database("ProjetoLTP2").Collection("Usuarios")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var results []bson.M
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func Upload(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(10 << 20) // Limite de 10MB
// 	if err != nil {
// 		http.Error(w, "Erro ao fazer o upload do arquivo", http.StatusBadRequest)
// 		return
// 	}

// 	file, handler, err := r.FormFile("image")
// 	if err != nil {
// 		http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Cria um diretório de uploads se não existir
// 	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
// 		err = os.Mkdir("uploads", os.ModePerm)
// 		if err != nil {
// 			http.Error(w, "Erro ao criar o diretório de uploads", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Cria um arquivo temporário no servidor para armazenar a imagem
// 	tempFile, err := os.CreateTemp("uploads", "upload-*.png")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Erro ao criar o arquivo temporário: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer tempFile.Close()

// 	// Copia o conteúdo do arquivo para o arquivo temporário
// 	_, err = io.Copy(tempFile, file)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Erro ao salvar o arquivo: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Cria a entrada para o banco de dados
// 	novaImagem := models.Imagem{
// 		Id:           primitive.NewObjectID(),
// 		FilePath:     tempFile.Name(),
// 		Img:          handler.Filename,
// 		Descricaoimg: r.FormValue("descricao"),
// 	}

// 	client := ConnectBd()
// 	collection := client.Database("ProjetoLTP2").Collection("Imagens")
// 	_, err = collection.InsertOne(context.Background(), novaImagem)
// 	if err != nil {
// 		http.Error(w, "Erro ao salvar a imagem no banco de dados", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode("Upload bem sucedido")
// }
