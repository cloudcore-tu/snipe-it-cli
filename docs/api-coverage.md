# API カバレッジ

Snipe-IT REST API v1 のコマンドカバレッジ一覧。
ルート定義は `routes/api.php`（github.com/snipe/snipe-it）を参照。

凡例: ✅ 実装済み / ⬜ 対象外（管理画面操作 / PDF / ファイルアップロード）

## CRUD リソース

| CLI コマンド | API パス | list | get | create | update | delete | 追加操作 |
|---|---|---|---|---|---|---|---|
| `assets` | `/api/v1/hardware` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout, checkin, audit, restore |
| `users` | `/api/v1/users` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `licenses` | `/api/v1/licenses` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout, checkin |
| `accessories` | `/api/v1/accessories` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout, checkin |
| `components` | `/api/v1/components` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout, checkin |
| `consumables` | `/api/v1/consumables` | ✅ | ✅ | ✅ | ✅ | ✅ | checkout |
| `categories` | `/api/v1/categories` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `companies` | `/api/v1/companies` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `locations` | `/api/v1/locations` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `manufacturers` | `/api/v1/manufacturers` | ✅ | ✅ | ✅ | ✅ | ✅ | restore |
| `models` | `/api/v1/models` | ✅ | ✅ | ✅ | ✅ | ✅ | restore |
| `departments` | `/api/v1/departments` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `statuslabels` | `/api/v1/statuslabels` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `suppliers` | `/api/v1/suppliers` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `fieldsets` | `/api/v1/fieldsets` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `maintenances` | `/api/v1/maintenances` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `fields` | `/api/v1/fields` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `depreciations` | `/api/v1/depreciations` | ✅ | ✅ | ✅ | ✅ | ✅ | |
| `groups` | `/api/v1/groups` | ✅ | ✅ | ✅ | ✅ | ✅ | |

## サブリソース

| CLI コマンド | API エンドポイント |
|---|---|
| `assets history --id N` | `GET /api/v1/hardware/{N}/history` |
| `assets licenses --id N` | `GET /api/v1/hardware/{N}/licenses` |
| `assets assigned-assets --id N` | `GET /api/v1/hardware/{N}/assigned/assets` |
| `assets assigned-accessories --id N` | `GET /api/v1/hardware/{N}/assigned/accessories` |
| `assets assigned-components --id N` | `GET /api/v1/hardware/{N}/assigned/components` |
| `assets bytag --tag TAG` | `GET /api/v1/hardware/bytag/{tag}` |
| `assets byserial --serial SERIAL` | `GET /api/v1/hardware/byserial/{serial}` |
| `users assets --id N` | `GET /api/v1/users/{N}/assets` |
| `users licenses --id N` | `GET /api/v1/users/{N}/licenses` |
| `users accessories --id N` | `GET /api/v1/users/{N}/accessories` |
| `users consumables --id N` | `GET /api/v1/users/{N}/consumables` |
| `licenses history --id N` | `GET /api/v1/licenses/{N}/history` |
| `licenses seats list --id N` | `GET /api/v1/licenses/{N}/seats` |
| `licenses seats get --id N --seat-id M` | `GET /api/v1/licenses/{N}/seats/{M}` |
| `licenses seats update --id N --seat-id M --data JSON` | `PATCH /api/v1/licenses/{N}/seats/{M}` |
| `accessories history --id N` | `GET /api/v1/accessories/{N}/history` |
| `accessories checkedout --id N` | `GET /api/v1/accessories/{N}/checkedout` |
| `components history --id N` | `GET /api/v1/components/{N}/history` |
| `components assets --id N` | `GET /api/v1/components/{N}/assets` |
| `consumables history --id N` | `GET /api/v1/consumables/{N}/history` |
| `consumables users --id N` | `GET /api/v1/consumables/{N}/users` |
| `locations users --id N` | `GET /api/v1/locations/{N}/users` |
| `locations assets --id N` | `GET /api/v1/locations/{N}/assets` |
| `locations assigned-assets --id N` | `GET /api/v1/locations/{N}/assigned/assets` |
| `locations assigned-accessories --id N` | `GET /api/v1/locations/{N}/assigned/accessories` |
| `locations history --id N` | `GET /api/v1/locations/{N}/history` |
| `statuslabels assetlist --id N` | `GET /api/v1/statuslabels/{N}/assetlist` |
| `statuslabels counts-by-label` | `GET /api/v1/statuslabels/assets/name` |
| `statuslabels counts-by-type` | `GET /api/v1/statuslabels/assets/type` |
| `fieldsets fields --id N` | `GET /api/v1/fieldsets/{N}/fields` |
| `maintenances history --id N` | `GET /api/v1/maintenances/{N}/history` |
| `models history --id N` | `GET /api/v1/models/{N}/history` |

## レポート・アカウント

| CLI コマンド | API エンドポイント |
|---|---|
| `reports activity` | `GET /api/v1/reports/activity` |
| `reports depreciation` | `GET /api/v1/reports/depreciation` |
| `account requestable` | `GET /api/v1/account/requestable/hardware` |
| `account requests` | `GET /api/v1/account/requests` |
| `account request --id N` | `POST /api/v1/account/request/{N}` |
| `account cancel-request --id N` | `POST /api/v1/account/request/{N}/cancel` |

## labels / imports / settings / notes

| CLI コマンド | API エンドポイント |
|---|---|
| `labels list` | `GET /api/v1/labels` |
| `labels get --name NAME [--output-file PATH]` | `GET /api/v1/labels/{name}` (PDF binary) |
| `imports list` | `GET /api/v1/imports` |
| `imports get --id N` | `GET /api/v1/imports/{N}` |
| `imports create --file PATH --type TYPE` | `POST /api/v1/imports` (multipart) |
| `imports update --id N --data JSON` | `PATCH /api/v1/imports/{N}` |
| `imports delete --id N --yes` | `DELETE /api/v1/imports/{N}` |
| `imports process --id N` | `POST /api/v1/imports/{N}/process` |
| `settings get` | `GET /api/v1/settings/general` |
| `settings update --data JSON` | `POST /api/v1/settings` |
| `settings login-attempts` | `GET /api/v1/settings/login-attempts` |
| `settings backups` | `GET /api/v1/settings/backups` |
| `settings backup-download [--name NAME] [--output-file PATH]` | `GET /api/v1/settings/backups/download/latest\|{name}` |
| `notes list --asset-id N` | `GET /api/v1/notes/{N}/index` |
| `notes create --asset-id N --data JSON` | `POST /api/v1/notes/{N}/store` |

## fields 追加操作

| CLI コマンド | API エンドポイント |
|---|---|
| `fields associate --id N --fieldset-id M` | `POST /api/v1/fields/{N}/associate` |
| `fields disassociate --id N --fieldset-id M` | `POST /api/v1/fields/{N}/disassociate` |
| `fields reorder --fieldset-id N --data '[1,2,3]'` | `POST /api/v1/fields/fieldsets/{N}/order` |

## account 追加操作

| CLI コマンド | API エンドポイント |
|---|---|
| `account eulas` | `GET /api/v1/account/eulas` |
| `account tokens` | `GET /api/v1/account/personal-access-tokens` |
| `account token-create --data JSON` | `POST /api/v1/account/personal-access-tokens` |
| `account token-delete --token-id N` | `DELETE /api/v1/account/personal-access-tokens/{N}` |
