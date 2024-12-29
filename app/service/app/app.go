package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/notnil/chess"
	"github.com/gorilla/websocket"
	"github.com/matoous/go-nanoid/v2"
)

type ChessGame struct {
	PlayerWhite     string
	PlayerBlack     string
	CurrentTurn     string
	Color           string
	InitialTime     int
	TimeControl     int
	GameStarted     bool
	PlayerWhiteTime time.Time
	PlayerBlackTime time.Time
	LastUpdate      time.Time

	Game			*chess.Game

	Connections []*websocket.Conn
	mu sync.Mutex
}


var GameStore = make(map[string]*ChessGame)

func NewApp() *ChessGame {
	return &ChessGame{}
}

func (c *ChessGame) CreateGame(initialTime, timeControl int, color string) error {
	gameID, err := gonanoid.New(10)
	if err != nil {
		return err
	}

	game := &ChessGame{
		InitialTime: initialTime,
		TimeControl: timeControl,
		Color:       color,
		GameStarted: false,
	}

	GameStore[gameID] = game

	fmt.Println(GameStore)
	return nil
}

func JoinGame(game *ChessGame, username string) error {

	
	if game.PlayerWhite != "" && game.PlayerBlack != "" {
		return fmt.Errorf("game already full")
	}

	assignPlayer(game, username)

	if game.PlayerWhite != "" && game.PlayerBlack != "" && !game.GameStarted {
		startGame(game)
	}
	fmt.Printf("%+v\n", game)
	return nil
}

func MakeMove(game *ChessGame, move string) (string, error) {

	chessGame := game.Game

	err := chessGame.MoveStr(move)
	if err != nil {
		return "", err
	}

	game.CurrentTurn = chessGame.Position().Turn().String()

	isGameOver(*game.Game)
 
	pos := chessGame.Position()
	return pos.String(), nil

}

func assignPlayer(game *ChessGame, username string) {
	if game.Color == "white" {
		if game.PlayerWhite == "" {
			game.PlayerWhite = username
		} else if game.PlayerBlack == "" && game.PlayerWhite != username {
			game.PlayerBlack = username
		}

	} else if game.Color == "black" {
		if game.PlayerBlack == "" {
			game.PlayerBlack = username
		} else if game.PlayerWhite == "" && game.PlayerBlack != username {
			game.PlayerWhite = username
		}
	}
}

func startGame(game *ChessGame) {

	nowTime := time.Now()
	playerTime := nowTime.Add(time.Duration(game.InitialTime) * time.Minute)

	fen := getFen()
	FEN, err := chess.FEN(fen)
	if err != nil {
		return
	}

	chessGame := chess.NewGame(FEN)
	turn := chessGame.Position().Turn().String()

	game.PlayerWhiteTime = playerTime
	game.PlayerBlackTime = playerTime
	game.LastUpdate = nowTime
	game.CurrentTurn = turn
	game.GameStarted = true

	game.Game = chessGame
}

func getFen() string {
	return "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
}

func isGameOver(game chess.Game) (string, error) {

	fmt.Println(game.Outcome())
	fmt.Println(game.Method())

	return "", nil
}
