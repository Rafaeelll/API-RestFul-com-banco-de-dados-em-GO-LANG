package main


// Instruções do trabalho da materia Topicos especiais: 
// O projeto deve apresentar um conjunto de APIs que manipulam dados em banco de dados usando a linguagem GO. 
// Pode-se utilizar qualquer dependência para criar os endpoints. Pode-se utilizar qualquer dependência para um ORM 
// que conecta e manipula dados em banco de dados.

// Passo a passo que utlizei para desenvolver o projeto:

// BD escolhido: Mysql hospedado no Clever Cloud.

// 1. Para conectar a API em Go a um banco de dados MySQL hospedado na Clever Cloud,
// precisei usar um driver de banco de dados MySQL e configurar a conexão, assim foi possivel interagir com o MySQL em Go.

// 2. Istalei o pacote com o seguinte comando: 
// go get -u github.com/go-sql-driver/mysql

// 3. Configurei a conexão com o banco de dados passando as credenciais que constam lá no BD do Clever Cloud, é possivel
// ver isso no codigo na linha 63.

// 4. A função initDB() garante que a conexão com o BD seja inicializada passando as credencias do BD dentro
// da função sql.Open(). 
// OBS: Chamei a função initDB() no início da função main() para garantir que a conexão com o banco de dados seja 
// estabelecida no início da execução do seu programa.

// 5. As demais funções garante que interagimos com o BD fazendo execuções como get, post, put e delete para isso
// utilizei querys do mysql como select, delete, insert e uptade.

// 6. Criei um endpoint para cada função que necessita executar querys no BD, como podemos ver no codigo a partir 
// da linha 306 foi criada as rotas com seus cujos metodos e funções.

// Tabela criadas no BD: 

// Criação da table usuario no banco de dados mysql (clever-cloud).
// CREATE TABLE usuario (
//   id int auto_increment not null,
//   nome varchar(30) not null,
//   sobrenome varchar(30) not null,
//   email varchar(200) unique not null,
//   senha varchar(200) not null,
//   telefone varchar(20) not null,
//   jogo_fav varchar(30) null,
//   plataforma_fav varchar(30) null,
//   PRIMARY KEY(id));

// Criação da table jogo no banco de dados mysql (clever-cloud).
// CREATE TABLE jogo (
//   id int auto_increment not null,
//   nome varchar(30) not null,
//   usuario_id int not null
//   PRIMARY KEY (id));

// Adicionado a foreign key na table jogo
// ALTER TABLE jogo add constraint fk_jogo_usuario_id foreign key (usuario_id) references usuario (id);

// Podemos dizer que um usuário pode ter vários jogos (relação "Um para Muitos"), mas cada jogo pertence a apenas um usuário.
// a intenção é modelar uma relação em que um usuário pode ter vários jogos associados a ele, 
// mas cada jogo está associado a apenas um usuário.

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//Importação do pacote driver mysql:
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Usuario representa dados sobre um usuario.
type Usuario struct {
    ID            int    `json:"id"`
    Nome          string `json:"nome"`
    Sobrenome     string `json:"sobrenome"`
    Email         string `json:"email"`
    Senha         string `json:"senha"`
    Telefone      string `json:"telefone"`
    PlataformaFav string `json:"plataforma_fav"`
    JogoFav       string `json:"jogo_fav"`
}
// Jogo representa dados sobre um jogo.

type Jogo struct{
	ID int `json:"id"`
	Nome string `json:"nome"`
	UsuarioID int `json:"usuario_id"`
}


// initDB configura a conexão com o banco de dados MySQL.
func initDB() {
	var err error
	db, err = sql.Open("mysql", "uhob3wlhzo5mjlef:lvaAI7hSG7iy5fEAg1ks@tcp(bpnlvjvv87p7vzryeliw-mysql.services.clever-cloud.com:3306)/bpnlvjvv87p7vzryeliw")
	if err != nil {
		log.Fatal(err)
	}

	// Testa a conexão com o banco de dados
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexão com o banco de dados bem-sucedIDa")
}

// getusuarios responde com a lista de todos os usuarios em JSON.
func getUsuarios(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM usuario") 
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var a Usuario
		err := rows.Scan(&a.ID,
			&a.Nome, 
			&a.Sobrenome, 
			&a.Email, 
			&a.Senha, 
			&a.Telefone,
			&a.PlataformaFav,
			&a.JogoFav)
		if err != nil {
			log.Fatal(err)
		}
		usuarios = append(usuarios, a)
	}

	c.IndentedJSON(http.StatusOK, usuarios)
}

// postusuarios adiciona um usuario a partir do JSON recebIDo no corpo da solicitação.
func postUsuarios(c *gin.Context) {
	var newusuarios Usuario

	// Chama BindJSON para vincular o JSON recebIDo a
	// newusuarios.
	if err := c.BindJSON(&newusuarios); err != nil {
		log.Fatal(err)
		return
	}

	// Adiciona o novo usuario à fatia.
	_, err := db.Exec("INSERT INTO usuario (id, nome, sobrenome, email, senha, telefone, plataforma_fav, jogo_fav ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		newusuarios.ID,
		newusuarios.Nome, 
		newusuarios.Sobrenome, 
		newusuarios.Email, 
		newusuarios.Senha, 
		newusuarios.Telefone,
		newusuarios.PlataformaFav, 
		newusuarios.JogoFav)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newusuarios)
}

// getusuariosByID localiza o usuario cujo valor de ID corresponde ao parâmetro de ID
// enviado pelo cliente, em seguIDa, retorna esse usuario como resposta.
func getUsuariosByID(c *gin.Context) {
	ID := c.Param("id")

	// Consulta o usuario no banco de dados pelo ID.
	var a Usuario
	err := db.QueryRow("SELECT * FROM usuario WHERE ID=?", ID).Scan(&a.ID, 
			&a.Nome, 
			&a.Sobrenome, 
			&a.Email, 
			&a.Senha, 
			&a.Telefone,
			&a.PlataformaFav,
			&a.JogoFav)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Usuario não encontrado"})
			return
		}
		log.Fatal(err)
		return
	}

	c.IndentedJSON(http.StatusOK, a)
} 

// putUsuarios atualiza um usuário existente com base no ID.
func putUsuarios(c *gin.Context) {
    var updatedUsuario Usuario

    if err := c.BindJSON(&updatedUsuario); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("ID do usuário a ser atualizado: %d\n", updatedUsuario.ID)

    _, err := db.Exec("UPDATE usuario SET nome=?, sobrenome=?, email=?, senha=?, telefone=?, plataforma_fav=?, jogo_fav=? WHERE ID=?",
        updatedUsuario.Nome,
        updatedUsuario.Sobrenome,
        updatedUsuario.Email,
        updatedUsuario.Senha,
        updatedUsuario.Telefone,
        updatedUsuario.PlataformaFav,
        updatedUsuario.JogoFav,
        updatedUsuario.ID)

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("Atualização bem-sucedIDa")

    c.IndentedJSON(http.StatusOK, updatedUsuario)
}


// deleteUsuarios exclui um usuário com base no ID.
func deleteUsuarios(c *gin.Context) {
    ID := c.Param("id")

    _, err := db.Exec("DELETE FROM usuario WHERE ID=?", ID)
    if err != nil {
        log.Fatal(err)
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Usuário excluído com sucesso"})
}

func getJogos(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM jogo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jogos []Jogo
	for rows.Next() {
		var a Jogo
		err := rows.Scan(&a.ID,
			&a.Nome, 
			&a.UsuarioID)
		if err != nil {
			log.Fatal(err)
		}
		jogos = append(jogos, a)
	}

	c.IndentedJSON(http.StatusOK, jogos)
}
func postJogos(c *gin.Context) {
	var newjogos Jogo

	// Chama BindJSON para vincular o JSON recebIDo a
	// newjogos.
	if err := c.BindJSON(&newjogos); err != nil {
		log.Fatal(err)
		return
	}

	// Adiciona o novo jogo à fatia.
	_, err := db.Exec("INSERT INTO jogo (id, nome, usuario_id) VALUES (?, ?, ?)",
		newjogos.ID,
		newjogos.Nome, 
		newjogos.UsuarioID)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newjogos)
}

func getJogosByID(c *gin.Context) {
	ID := c.Param("id")

	// Consulta o jogo no banco de dados pelo ID.
	var a Jogo
	err := db.QueryRow("SELECT * FROM jogo WHERE ID=?", ID).Scan(&a.ID, 
			&a.Nome, 
			&a.UsuarioID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Usuario não encontrado"})
			return
		}
		log.Fatal(err)
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func putJogos(c *gin.Context) {

    var updatedJogo Jogo

    if err := c.BindJSON(&updatedJogo); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	fmt.Printf("ID do usuário a ser atualizado: %d\n", updatedJogo.ID)

    _, err := db.Exec("UPDATE jogo SET nome=?, usuario_id=? WHERE ID=?",
		updatedJogo.Nome,
		updatedJogo.UsuarioID,
		updatedJogo.ID)

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	fmt.Println("Atualização bem-sucedIDa")

    c.IndentedJSON(http.StatusOK, updatedJogo)
}

func deleteJogos(c *gin.Context) {
    ID := c.Param("id")

    _, err := db.Exec("DELETE FROM jogo WHERE ID=?", ID)
    if err != nil {
        log.Fatal(err)
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Jogo excluído com sucesso"})
}
func main() {
	initDB()

	router := gin.Default()
	router.GET("/usuarios", getUsuarios)
	router.GET("/usuarios/:id", getUsuariosByID)
	router.POST("/usuarios", postUsuarios)
	router.PUT("/usuarios/:id", putUsuarios)
	router.DELETE("/usuarios/:id", deleteUsuarios)
	router.GET("/jogos", getJogos)
	router.GET("/jogos/:id", getJogosByID)
	router.POST("/jogos", postJogos)
	router.PUT("/jogos/:id", putJogos)
	router.DELETE("/jogos/:id", deleteJogos)
	router.Run("localhost:8080")
}
