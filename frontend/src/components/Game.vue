<template>
  <div class="m-8">
    <p class="text-2xl font-black text-center">縱橫殺棋對抗賽</p>
    <div class="mt-8 flex justify-around">
      <div class="w-1/3">
        <p class="text-xl text-center">輸入棋盤狀態</p>
        <el-input v-model="data" :rows="8" type="textarea" class="mt-2" />
        <el-button type="primary" plain round class="mt-2 w-full" @click="initBoard">初始化棋盤</el-button>
        <div class="mt-8 w-full flex space-x-8 align-middle">
          <el-input-number v-model="randBoardH" :min="1" :max="8" />
          <div class="m-auto">X</div>
          <el-input-number v-model="randBoardW" :min="1" :max="8" />
        </div>
        <el-button type="info" plain round class="mt-2 !ml-0 w-full" @click="randBoard">隨機棋盤</el-button>
        <el-button type="warning" plain round class="mt-8 !ml-0 w-full" @click="aiSelect">AI</el-button>
        <div v-if="aiElapsedTime !== ''" class="mt-2">
          <el-alert type="info" :closable="false" center>
            <p class="text-lg">Points : {{ aiBestScore }}</p>
            <p class="text-lg">Seconds: {{ aiElapsedTime }}</p>
          </el-alert>
        </div>
      </div>
      <div class="w-2/3 mx-8">
        <div v-if="n > 0 && m > 0" class="w-full">
          <div class="text-center">
            <p class="text-lg">
              輪到
              <span v-if="turn">先手</span>
              <span v-else>後手</span>
            </p>
          </div>
          <table class="w-full mt-4 border-collapse border border-slate-400">
            <tr v-for="i in n + 1">
              <td v-for="j in m + 1">
                <div class="w-ful h-full border border-slate-300">
                  <div class="w-full h-8 flex align-middle text-center">
                    <div v-if="i === 1 || j === 1" class="m-auto text-center w-full">
                      <el-button v-if="i === 1 && j > 1" class="w-full h-full" @click="kill(i, j)">
                        v
                      </el-button>
                      <el-button v-if="j === 1 && i > 1" class="w-full h-full" @click="kill(i, j)">
                        >
                      </el-button>
                    </div>
                    <div v-else class="m-auto text-center">
                      <div v-if="board[i - 2][j - 2] == 1" class="w-4 h-4 bg-black rounded-full"> </div>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </table>
          <div class="mt-4 text-center">
            <div class="text-lg">
              比分： {{ firstScore }} ： {{ secondScore }}
            </div>
          </div>
          <div v-if="isGameEnd" class="mt-4">
            <el-alert type="warning" :closable="false" center>
              <p class="text-lg">
                遊戲結束
                <span v-if="firstScore > secondScore">先手勝</span>
                <span v-else-if="secondScore > firstScore">後手勝</span>
                <span v-else>和局</span>
              </p>
            </el-alert>
          </div>
        </div>
        <div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { AISelect } from '../../wailsjs/go/main/App'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const data = ref("")
const board = ref<number[][]>([])
const n = ref(0)
const m = ref(0)
const turn = ref(true)
const firstScore = ref(0)
const secondScore = ref(0)
const isGameEnd = ref(false)

const initBoard = () => {
  data.value = data.value.trim()
  if (data.value === "") {
    ElMessage({
      message: "請輸入棋盤狀態",
      type: "error",
      duration: 1000,
    })
    return
  }
  board.value = []
  board.value = data.value.split("\n").map(line => line.split(" ").map(x => parseInt(x)))
  n.value = board.value.length
  m.value = board.value[0].length
  turn.value = true
  firstScore.value = 0
  secondScore.value = 0
  isGameEnd.value = false
  ElMessage({
    message: "初始化棋盤成功",
    type: "success",
    duration: 1000,
  })
  aiBestScore.value = 0
  aiElapsedTime.value = ""
}

const randBoardW = ref(4)
const randBoardH = ref(4)

const randBoard = () => {
  data.value = ""
  for (let i = 0; i < randBoardH.value; i++) {
    if (i !== 0) data.value += "\n"
    for (let j = 0; j < randBoardW.value; j++) {
      if (j !== 0) data.value += " "
      data.value += Math.floor(Math.random() * 2)
    }
  }
  if (!data.value.includes("1")) randBoard()
  initBoard()
}

const killCheck = (i: number, j: number): boolean => {
  if (j === 1) {
    for (let k = 0; k < m.value; k++) {
      if (board.value[i - 2][k] === 1) return true
    }
  } else if (i === 1) {
    for (let k = 0; k < n.value; k++) {
      if (board.value[k][j - 2] === 1) return true
    }
  }
  return false
}

const gameEndCheck = (): boolean => {
  return board.value.every(x => x.every(y => y === 0))
}

const kill = (i: number, j: number) => {
  if (!killCheck(i, j)) {
    ElMessage({
      message: "無法殺棋",
      type: "error",
      duration: 1000,
    })
    return
  }
  if (j === 1) {
    for (let k = 0; k < m.value; k++) {
      if (board.value[i - 2][k] === 1) {
        if (turn.value) firstScore.value++
        else secondScore.value++
        board.value[i - 2][k] = 0
      }
    }
  } else if (i === 1) {
    for (let k = 0; k < n.value; k++) {
      if (board.value[k][j - 2] === 1) {
        if (turn.value) firstScore.value++
        else secondScore.value++
        board.value[k][j - 2] = 0
      }
    }
  }
  if (gameEndCheck()) {
    isGameEnd.value = true
    return
  }
  turn.value = !turn.value
}

const aiBestScore = ref(0)
const aiElapsedTime = ref("")

const aiSelect = () => {
  AISelect(board.value, firstScore.value, secondScore.value)
    .then((res) => {
      if (res === null) return
      if (res["bestMove"].includes("Row")) kill(parseInt(res["bestMove"].split(" ")[1]) + 1, 1)
      else kill(1, parseInt(res["bestMove"].split(" ")[1]) + 1)
      aiBestScore.value = res["bestScore"]
      aiElapsedTime.value = res["elapsedTime"]
    })
}
</script>
