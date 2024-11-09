package main
import (
	"math/rand"
	"time"
	
)

func randomWordGenerate() string {
    words := []string{"下北沢", "ヘッドフォン", "データベース", "マックブック","スマートフォン","ノートパソコン","ワイヤレスイヤホン","ワイヤレスマウス","ワイヤレスキーボード"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}