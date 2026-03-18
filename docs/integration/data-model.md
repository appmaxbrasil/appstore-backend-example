# Data Model

Database models, field relationships, and order status lifecycle.

Models: `app/models/`

---

## Entity Relationship Diagram

```
+--------------------+         +--------------------+         +---------------------+
|   Installation     |         |      Order         |         |   WebhookEvent      |
+--------------------+         +--------------------+         +---------------------+
| id            PK   |<---+    | id            PK   |         | id             PK   |
| external_key  UQ   |    +----| installation_id FK |         | event               |
| app_id             |         | appmax_customer_id |         | event_type          |
| merchant_client_id |         | appmax_order_id UQ |         | appmax_order_id     |
| merchant_client_   |         | status             |         | payload      JSONB  |
|   secret           |         | payment_method     |         | processed           |
| external_id   UQ   |         | total_cents        |         | processed_at        |
| installed_at       |         | pix_qr_code        |         | error_message       |
| created_at         |         | pix_emv            |         | created_at          |
| updated_at         |         | boleto_pdf_url     |         +---------------------+
+--------------------+         | boleto_digitavel   |
                               | upsell_hash        |
                               | created_at         |
                               | updated_at         |
                               +--------------------+
```

---

## Installation

Represents a merchant's app installation. Created during the install flow (via either
the browser OAuth path or the Appmax health check POST).

| Column                 | Type      | Constraints               | Description                                    |
|------------------------|-----------|---------------------------|------------------------------------------------|
| `id`                   | int64     | PK, auto-increment        | Internal ID                                    |
| `external_key`         | string    | NOT NULL, UNIQUE INDEX    | Developer-provided merchant identifier          |
| `app_id`               | string    | NOT NULL                  | Appmax app ID (UUID or numeric, depending on creation path) |
| `merchant_client_id`   | string    | NOT NULL                  | OAuth client ID for merchant API calls          |
| `merchant_client_secret` | string  | NOT NULL                  | OAuth client secret for merchant API calls      |
| `external_id`          | string    | NOT NULL, UNIQUE INDEX    | UUID generated on first create, returned to Appmax |
| `installed_at`         | timestamp | NOT NULL                  | When installation was created/last updated      |
| `created_at`           | timestamp | NOT NULL                  | Record creation time                            |
| `updated_at`           | timestamp | NOT NULL                  | Record last update time                         |

**Lookup key**: `external_key` (used by `MerchantContext` middleware and all `{key}` routes).

**Upsert behavior**: On duplicate `external_key`, updates `merchant_client_id`,
`merchant_client_secret`, and `installed_at`. Does not regenerate `external_id`.

Model: `app/models/installation.go`

---

## Order

Represents an Appmax checkout order. Created best-effort after payment processing.

| Column               | Type      | Constraints               | Description                                    |
|----------------------|-----------|---------------------------|------------------------------------------------|
| `id`                 | int64     | PK, auto-increment        | Internal ID                                    |
| `installation_id`    | int64     | NOT NULL, FK              | References `installations.id`                  |
| `appmax_customer_id` | int       | NOT NULL                  | Customer ID in Appmax system                   |
| `appmax_order_id`    | int       | NOT NULL, UNIQUE INDEX    | Order ID in Appmax system (deduplication key)  |
| `status`             | string    | NOT NULL, default `pendente` | Current order status                        |
| `payment_method`     | string    | nullable                  | `credit_card`, `pix`, or `boleto`              |
| `total_cents`        | int       | NOT NULL, default 0       | Order amount in cents                          |
| `pix_qr_code`        | string    | nullable                  | Base64-encoded QR code image (Pix only)        |
| `pix_emv`            | string    | nullable                  | Copy-paste Pix code (Pix only)                 |
| `boleto_pdf_url`     | string    | nullable                  | URL to boleto PDF (Boleto only)                |
| `boleto_digitavel`   | string    | nullable                  | Boleto barcode line (Boleto only)              |
| `upsell_hash`        | string    | nullable                  | Hash for creating upsell offers                |
| `created_at`         | timestamp | NOT NULL                  | Record creation time                           |
| `updated_at`         | timestamp | NOT NULL                  | Record last update time                        |

**Deduplication key**: `appmax_order_id` (unique index, used for webhook event matching).

**Best-effort persistence**: If the DB write fails during payment, the payment response is
still returned to the client (warning logged).

**Status updates**: Updated by webhook processing when matching events arrive.

Model: `app/models/order.go`

---

## WebhookEvent

Audit trail for all incoming webhook events from Appmax.

| Column            | Type       | Constraints               | Description                                     |
|-------------------|------------|---------------------------|-------------------------------------------------|
| `id`              | int64      | PK, auto-increment        | Internal ID                                     |
| `event`           | string     | NOT NULL                  | Event name (e.g., `OrderPaid`, `order_paid`)    |
| `event_type`      | string     | NOT NULL                  | Event category (`order`, `customer`, etc.)       |
| `appmax_order_id` | *int       | nullable                  | Associated Appmax order ID (null for non-order events) |
| `payload`         | JSONB      | NOT NULL                  | Full webhook payload as received                 |
| `processed`       | bool       | NOT NULL, default false   | Whether the event has been processed             |
| `processed_at`    | *timestamp | nullable                  | When processing completed                        |
| `error_message`   | string     | nullable                  | Error details if processing failed               |
| `created_at`      | timestamp  | NOT NULL                  | Record creation time                             |

**Deduplication query**: `WHERE event = ? AND appmax_order_id = ? AND processed = true AND id != ?`

**Processing states**:
- `processed = false`: Event received but not yet handled
- `processed = true, error_message = ""`: Successfully processed
- `processed = true, error_message != ""`: Processed with error (e.g., order update failed)

Model: `app/models/webhook_event.go`

---

## Key Field Relationships

| Lookup Pattern                                    | Purpose                                      |
|---------------------------------------------------|----------------------------------------------|
| `Installation.external_key` = route `{key}`       | Middleware loads installation for all protected routes |
| `Order.appmax_order_id` = webhook `data.order_id` | Match incoming webhooks to local orders       |
| `Order.installation_id` = `Installation.id`       | Scope orders to a specific merchant           |
| `WebhookEvent.appmax_order_id` = `Order.appmax_order_id` | Link webhook audit trail to orders  |

---

## Order Status Lifecycle

```
                    +------------+
                    |  pendente  |
                    +-----+------+
                          |
            +-------------+-------------+
            |             |             |
     OrderBillet    OrderPix      OrderPaid/
     Created       Created       OrderApproved/
            |             |       CreditCard
            |             |             |
            v             v             v
     +------+------+  +--+---+   +-----+------+
     | pendente    |  |pend. |   | aprovado   |
     | (boleto)    |  |(pix) |   |            |
     +------+------+  +--+---+   +-----+------+
            |             |             |
     +------+------+  +--+---+   +-----+------+
     |OrderBillet   |  |Order |   |OrderPending|
     |Overdue       |  |Pix   |   |Integration |
     |              |  |Expired|   |            |
     v              |  v      |   v            |
  cancelado         | cancel. |  pendente_     |
                    |         |  integracao    |
                    +---------+        |       |
                                       v       |
                              +--------+--+    |
                              |integrado  |    |
                              +-----------+    |
                                               |
                                    +----------+
                                    |
                              +-----v------+
                              |OrderRefund |
                              |            |
                              v            |
                          estornado        |
                                           |
                              +------------+
                              |OrderCharge |
                              |BackIn      |
                              |Treatment   |
                              v
                         chargeback_
                         em_tratativa
```

### Status Descriptions

| Status                    | Description                                               | Terminal? |
|---------------------------|-----------------------------------------------------------|-----------|
| `pendente`                | Awaiting payment (Pix QR code generated, boleto issued)   | No        |
| `aprovado`                | Payment approved, funds available                         | No        |
| `autorizado`              | Credit card authorized, anti-fraud analysis in progress   | No        |
| `cancelado`               | Declined, not authorized, expired (Pix/Boleto)            | Yes       |
| `estornado`               | Refunded (full or partial)                                | Yes       |
| `integrado`               | Order fully integrated, ready to ship                     | Yes       |
| `pendente_integracao`     | Approved but pending integration issues                   | No        |
| `chargeback_em_tratativa` | Chargeback under analysis by Appmax                       | No        |
