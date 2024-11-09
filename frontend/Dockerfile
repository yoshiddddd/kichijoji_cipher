# Node.jsの公式イメージをベースに使用
FROM node:18

# 作業ディレクトリを設定
WORKDIR /app

# パッケージファイルをコピーして依存関係をインストール
COPY package*.json ./

# 依存関係をインストール
RUN npm install

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションがリッスンするポートを指定（必要に応じて）
EXPOSE 3000

# アプリケーションを起動
CMD ["npm", "run", "dev"]
