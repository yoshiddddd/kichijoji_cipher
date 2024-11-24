package main
import (
	"math/rand"
	"time"
	
)

func firstRandomWordGenerate() string {
    words := []string{"ヘッドフォン","サイヤ人","なあぜなあぜ","イカサマ","テニス","魚","格差", "身長", "ドラム", "嘘つき", "東京", "逆上がり"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func secondRandomWordGenerate() string {
    words := []string{"ノートパソコン","俺が基準","有頂天","カメラ撮影","当たり前","燃えるごみ","経験値","万歩計", "インド人", "トップランク"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func thirdRandomWordGenerate() string {
	words := []string{"ワイヤレスイヤホン","LINEスタンプ","ワイヤレスキーボード","最強の要塞","モノマネ芸人","ターンテーブル","きんかんのど飴","ショーケース", "ローキック"}
    
    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}