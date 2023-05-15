package main

import (
	"context"
	"fmt"
	"math"
	"time"
)

type GameState struct {
	board       [][]int
	turn        bool
	firstScore  int
	secondScore int
	depth       int
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

var bestMove = -1

type AISelectResult struct {
	BestMove    string `json:"bestMove"`
	BestScore   int    `json:"bestScore"`
	ElapsedTime string `json:"elapsedTime"`
}

func (a *App) AISelect(board [][]int, firstScore, secondScore int) AISelectResult {
	bestMove = -1
	initState := &GameState{
		board:       board,
		turn:        true,
		firstScore:  firstScore,
		secondScore: secondScore,
		depth:       0,
	}
	startTime := time.Now()
	bestScore := alphaBeta(initState, math.MinInt, math.MaxInt)
	elapseTime := time.Since(startTime)

	return AISelectResult{
		BestMove:    printBestMove(len(initState.board)),
		BestScore:   bestScore,
		ElapsedTime: fmt.Sprintf("%f", elapseTime.Seconds()),
	}
}

func printBestMove(rows int) string {
	if bestMove == -1 {
		return "No move"
	}
	if bestMove < rows {
		return fmt.Sprintf("Row %d", bestMove+1)
	} else {
		return fmt.Sprintf("Column %d", bestMove-rows+1)
	}
}

// 判斷遊戲是否結束
func (s *GameState) isTerminal() bool {
	for i := 0; i < len(s.board); i++ {
		for j := 0; j < len(s.board[0]); j++ {
			if s.board[i][j] == 1 {
				return false
			}
		}
	}
	return true
}

func (s *GameState) getScore(n int) int {
	cnt := 0
	if n < len(s.board) { // row
		for i := 0; i < len(s.board[0]); i++ {
			if s.board[n][i] == 1 {
				cnt++
			}
		}
	} else { // col
		for i := 0; i < len(s.board); i++ {
			if s.board[i][n-len(s.board)] == 1 {
				cnt++
			}
		}
	}
	return cnt
}

func (s *GameState) nextBoard(n int) [][]int {
	newBoard := make([][]int, len(s.board))
	for i := range s.board {
		newBoard[i] = make([]int, len(s.board[i]))
		copy(newBoard[i], s.board[i])
	}
	if n < len(s.board) { // row
		for i := 0; i < len(s.board[0]); i++ {
			newBoard[n][i] = 0
		}
	} else { // col
		for i := 0; i < len(s.board); i++ {
			newBoard[i][n-len(s.board)] = 0
		}
	}
	return newBoard
}

func alphaBeta(state *GameState, alpha, beta int) int {
	if state.isTerminal() {
		return state.firstScore - state.secondScore
	}
	if state.turn { // 先手
		maxVal := math.MinInt
		for i := 0; i < len(state.board)+len(state.board[0]); i++ {
			cnt := state.getScore(i)
			if cnt == 0 {
				continue
			}
			newState := &GameState{
				board:       state.nextBoard(i),
				turn:        !state.turn,
				firstScore:  state.firstScore + cnt,
				secondScore: state.secondScore,
				depth:       state.depth + 1,
			}
			score := alphaBeta(newState, alpha, beta)
			if score > maxVal {
				maxVal = score
			}
			if maxVal > alpha {
				alpha = maxVal
				if state.depth == 0 {
					bestMove = i
				}
			}
			if alpha >= beta {
				break
			}
		}
		return maxVal
	} else { // 後手
		minVal := math.MaxInt
		for i := 0; i < len(state.board)+len(state.board[0]); i++ {
			cnt := state.getScore(i)
			if cnt == 0 {
				continue
			}
			newState := &GameState{
				board:       state.nextBoard(i),
				turn:        !state.turn,
				firstScore:  state.firstScore,
				secondScore: state.secondScore + cnt,
				depth:       state.depth + 1,
			}
			score := alphaBeta(newState, alpha, beta)
			if score < minVal {
				minVal = score
			}
			if minVal < beta {
				beta = minVal
			}
			if alpha >= beta {
				break
			}
		}
		return minVal
	}
}
