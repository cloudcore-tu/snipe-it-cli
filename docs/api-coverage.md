# API カバレッジ

Snipe-IT REST API v1 のコマンドカバレッジ一覧。

凡例: ✅ 実装済み / ⬜ 未実装

## リソース別カバレッジ

| CLI コマンド | API パス | list | get | create | update | delete | 備考 |
|---|---|---|---|---|---|---|---|
| `assets` | `/api/v1/hardware` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout/checkin/audit アクション付き |
| `users` | `/api/v1/users` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `licenses` | `/api/v1/licenses` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout/checkin アクション付き |
| `accessories` | `/api/v1/accessories` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout/checkin アクション付き |
| `components` | `/api/v1/components` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout/checkin アクション付き |
| `consumables` | `/api/v1/consumables` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `categories` | `/api/v1/categories` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `companies` | `/api/v1/companies` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `locations` | `/api/v1/locations` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `manufacturers` | `/api/v1/manufacturers` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `models` | `/api/v1/models` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `departments` | `/api/v1/departments` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `statuslabels` | `/api/v1/statuslabels` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `suppliers` | `/api/v1/suppliers` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `fieldsets` | `/api/v1/fieldsets` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `maintenances` | `/api/v1/maintenances` | ✅ | ✅ | ✅ | ✅ | ✅ | |

## 未実装エンドポイント

| エンドポイント | 説明 |
|---|---|
| `GET /api/v1/account/requestable/hardware` | リクエスト可能な資産一覧 |
| `GET /api/v1/hardware/bytag/{tag}` | 資産タグで資産取得 |
| `GET /api/v1/hardware/byserial/{serial}` | シリアル番号で資産取得 |
| `GET /api/v1/users/{id}/assets` | ユーザーの割り当て資産一覧 |
| `GET /api/v1/users/{id}/licenses` | ユーザーの割り当てライセンス一覧 |
| `GET /api/v1/fields` | カスタムフィールド一覧 |
| `GET /api/v1/reports/activity` | アクティビティレポート |
