# オニオンアーキテクチャを用いたブログアプリのAPI

## 概要

オニオンアーキテクチャを採用したブログアプリのAPI。

## アーキテクチャ

1. **ドメイン層 (Domain Layer)**
   - アプリケーションの中心となるビジネスロジックを含む
   - エンティティ、値オブジェクト、ドメインサービス、リポジトリで構成
   - 外部の層に依存しない
   - リポジトリ
     - ドメインオブジェクトの永続化を担当
     - インターフェースのみを定義

2. **ユースケース層 (Usecase Layer)**
   - アプリケーションの機能を実現するためのビジネスロジックを含む
   - ドメイン層とリポジトリ層に依存する

3. **インフラストラクチャ層 (Infrastructure Layer)**
   - データベース、外部APIなどの技術的な実装を含む
   - リポジトリインターフェースの実装を提供

4. **UI層 (UI Layer)**
   - ユーザーインターフェースを提供
   - HTTPハンドラ、ミドルウェアなどを含む

## API エンドポイント

### ユーザー関連

- `POST /api/users/register` - ユーザー登録
- `POST /api/users/login` - ログイン
- `GET /api/users/:id` - ユーザー情報取得
- `PUT /api/users/:id` - ユーザー情報更新
- `DELETE /api/users/:id` - ユーザー削除

### ブログ関連

- `POST /api/blogs` - ブログ投稿
- `GET /api/blogs` - ブログ一覧取得
- `GET /api/blogs/:id` - ブログ詳細取得
- `GET /api/users/:id/blogs` - ユーザーのブログ一覧取得
- `PUT /api/blogs/:id` - ブログ更新
- `DELETE /api/blogs/:id` - ブログ削除

### コメント関連

- `POST /api/blogs/:id/comments` - コメント投稿
- `GET /api/blogs/:id/comments` - ブログのコメント一覧取得
- `PUT /api/comments/:id` - コメント更新
- `DELETE /api/comments/:id` - コメント削除
