package main
import (
	"math/rand"
	"time"
	
)

func firstRandomWordGenerate() string {
    words := []string{"下北沢", "ヘッドフォン"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func secondRandomWordGenerate() string {
    words := []string{"ノートパソコン"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func thirdRandomWordGenerate() string {
	words := []string{"ワイヤレスイヤホン","ワイヤレスマウス","ワイヤレスキーボード"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}