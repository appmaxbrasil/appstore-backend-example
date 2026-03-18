# Appmax Webhooks — Event Reference

Appmax sends webhooks to `POST /webhooks/appmax`. The JSON payload structure depends on which **Content Model** is configured in the Appmax dashboard (Settings → Webhooks → Content Model). There are 4 configurable models plus one pre-configurable legacy format.

---

## Overview of Payload Models

| Model | Event format | `event_type` | Order ID field | Customer data |
|-------|-------------|-------------|----------------|--------------|
| [Standard](#1-standard) | PascalCase | `""` | `data.id` | Nested `data.customer` object |
| [Standard with Meta](#2-standard-with-meta) | PascalCase | `""` | `data.id` | Nested `data.customer` object + `data.meta: []` |
| [Two-Level Flat](#3-two-level-flat) | PascalCase | `""` | `data.order_id` | Flat `data.customer_*` + `data.customer_site` |
| [Custom Content](#4-custom-content) | PascalCase | `""` | `data.order_id` | Absent or minimal |
| [Old Legacy](#5-old-legacy) | snake_case | `"order"` | `data.order_id` | Absent |

**Detection logic (priority order):**
1. `event_type == "order"` → Old Legacy
2. `data.id` present AND `data.customer_id` present AND `data.meta` present → Standard with Meta
3. `data.id` present AND `data.customer_id` present → Standard
4. `data.order_id` present AND `data.order_total_products` present → Two-Level Flat
5. `data.order_id` present (no `order_*` fields) → Custom Content

**Order ID extraction rule:**
- `data.order_id` is preferred (Two-Level Flat, Custom Content, Old Legacy)
- `data.id` is used as order ID only when `data.customer_id` is also present (Standard, Standard with Meta) — this guards against treating customer events (where `data.id` is the customer ID) as order events

---

## 1. Standard

**Detection:** `data.id` present + `data.customer_id` present, no `data.meta`

**Envelope:**
```json
{
  "event": "OrderApproved",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

### Order Events

**OrderApproved**
```json
{
  "event": "OrderApproved",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderAuthorized**
```json
{
  "event": "OrderAuthorized",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "autorizado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderPaid**
```json
{
  "event": "OrderPaid",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderBilletCreated**
```json
{
  "event": "OrderBilletCreated",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "pendente",
    "payment_type": "Billet",
    "billet_url": "https://boleto.example.com/abc123",
    "billet_expiration_date": "2026-03-24",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderBilletOverdue**
```json
{
  "event": "OrderBilletOverdue",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "cancelado",
    "payment_type": "Billet",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderPixCreated**
```json
{
  "event": "OrderPixCreated",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "pendente",
    "payment_type": "Pix",
    "pix_code": "00020126330014br.gov.bcb.pix0111999999999952040000530398654041.005802BR5913Leandro Silva6008Sao Paulo62070503***6304ABCD",
    "pix_expiration_date": "2026-03-17T23:59:59",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderPaidByPix**
```json
{
  "event": "OrderPaidByPix",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "Pix",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderPixExpired**
```json
{
  "event": "OrderPixExpired",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "cancelado",
    "payment_type": "Pix",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderPendingIntegration**
```json
{
  "event": "OrderPendingIntegration",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "pendente_integracao",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderIntegrated**
```json
{
  "event": "OrderIntegrated",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "integrado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderRefund**
```json
{
  "event": "OrderRefund",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "estornado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderChargeBackInTreatment**
```json
{
  "event": "OrderChargeBackInTreatment",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "chargeback_em_tratativa",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**OrderUpSold**
```json
{
  "event": "OrderUpSold",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

### Customer Events

Customer events use `data.id` as the customer ID (not order ID). `ExtractOrderID` returns nil for these.

**CustomerCreated**
```json
{
  "event": "CustomerCreated",
  "event_type": "",
  "data": {
    "id": 7,
    "site_id": 1470,
    "firstname": "Leandro",
    "lastname": "Silva",
    "email": "leandro@example.com",
    "phone": "11999999999"
  }
}
```

**CustomerInterested**
```json
{
  "event": "CustomerInterested",
  "event_type": "",
  "data": {
    "id": 7,
    "site_id": 1470,
    "firstname": "Leandro",
    "lastname": "Silva",
    "email": "leandro@example.com",
    "phone": "11999999999"
  }
}
```

**CustomerContacted**
```json
{
  "event": "CustomerContacted",
  "event_type": "",
  "data": {
    "id": 7,
    "site_id": 1470,
    "firstname": "Leandro",
    "lastname": "Silva",
    "email": "leandro@example.com",
    "phone": "11999999999"
  }
}
```

### Subscription Events

**CreatedSubscription**
```json
{
  "event": "CreatedSubscription",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    }
  }
}
```

**SubscriptionCancellationEvent** (no-op — uses customer structure, no order ID)
```json
{
  "event": "SubscriptionCancellationEvent",
  "event_type": "",
  "data": {
    "id": 7,
    "site_id": 1470,
    "firstname": "Noeli",
    "lastname": "Guerra",
    "subscription": {
      "id": null
    }
  }
}
```

**SubscriptionDelayedEvent** (no-op)
```json
{
  "event": "SubscriptionDelayedEvent",
  "event_type": "",
  "data": {
    "id": 7,
    "site_id": 1470,
    "firstname": "Noeli",
    "lastname": "Guerra",
    "subscription": {
      "id": 99
    }
  }
}
```

---

## 2. Standard with Meta

**Detection:** `data.id` present + `data.customer_id` present + `data.meta` present

Same as Standard but with `"meta": []` appended to `data`. Customer events do not appear in this model.

**OrderApproved**
```json
{
  "event": "OrderApproved",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    },
    "meta": []
  }
}
```

**OrderBilletCreated**
```json
{
  "event": "OrderBilletCreated",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "pendente",
    "payment_type": "Billet",
    "billet_url": "https://boleto.example.com/abc123",
    "billet_expiration_date": "2026-03-24",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    },
    "meta": []
  }
}
```

**OrderPixCreated**
```json
{
  "event": "OrderPixCreated",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "pendente",
    "payment_type": "Pix",
    "pix_code": "00020126330014br.gov.bcb.pix0111999999999952040000530398654041.005802BR5913Leandro Silva6008Sao Paulo62070503***6304ABCD",
    "pix_expiration_date": "2026-03-17T23:59:59",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    },
    "meta": []
  }
}
```

**CreatedSubscription**
```json
{
  "event": "CreatedSubscription",
  "event_type": "",
  "data": {
    "id": 12844,
    "customer_id": 7,
    "total_products": 265.00,
    "total": 267.48,
    "status": "aprovado",
    "payment_type": "CreditCard",
    "customer": {
      "id": 7,
      "firstname": "Leandro",
      "lastname": "Silva",
      "email": "leandro@example.com",
      "phone": "11999999999"
    },
    "meta": []
  }
}
```

---

## 3. Two-Level Flat

**Detection:** `data.order_id` present + `data.order_total_products` present

All order fields use the `order_*` prefix. Customer fields use the `customer_*` prefix. No nested objects except `customer_site`.

**OrderApproved**
```json
{
  "event": "OrderApproved",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_total_products": 265.00,
    "order_total": 267.48,
    "order_status": "aprovado",
    "order_payment_type": "CreditCard",
    "customer_id": 7,
    "customer_firstname": "Leandro",
    "customer_lastname": "Silva",
    "customer_email": "leandro@example.com",
    "customer_phone": "11999999999",
    "customer_site": {
      "id": 1470
    }
  }
}
```

**OrderBilletCreated**
```json
{
  "event": "OrderBilletCreated",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_total_products": 265.00,
    "order_total": 267.48,
    "order_status": "pendente",
    "order_payment_type": "Billet",
    "order_billet_url": "https://boleto.example.com/abc123",
    "order_billet_expiration_date": "2026-03-24",
    "customer_id": 7,
    "customer_firstname": "Leandro",
    "customer_lastname": "Silva",
    "customer_email": "leandro@example.com",
    "customer_phone": "11999999999",
    "customer_site": {
      "id": 1470
    }
  }
}
```

**OrderPixCreated**
```json
{
  "event": "OrderPixCreated",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_total_products": 265.00,
    "order_total": 267.48,
    "order_status": "pendente",
    "order_payment_type": "Pix",
    "order_pix_code": "00020126330014br.gov.bcb.pix0111999999999952040000530398654041.005802BR5913Leandro Silva6008Sao Paulo62070503***6304ABCD",
    "order_pix_expiration_date": "2026-03-17T23:59:59",
    "customer_id": 7,
    "customer_firstname": "Leandro",
    "customer_lastname": "Silva",
    "customer_email": "leandro@example.com",
    "customer_phone": "11999999999",
    "customer_site": {
      "id": 1470
    }
  }
}
```

**OrderRefund**
```json
{
  "event": "OrderRefund",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_total_products": 265.00,
    "order_total": 267.48,
    "order_status": "estornado",
    "order_payment_type": "CreditCard",
    "customer_id": 7,
    "customer_firstname": "Leandro",
    "customer_lastname": "Silva",
    "customer_email": "leandro@example.com",
    "customer_phone": "11999999999",
    "customer_site": {
      "id": 1470
    }
  }
}
```

**CreatedSubscription**
```json
{
  "event": "CreatedSubscription",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_total_products": 265.00,
    "order_total": 267.48,
    "order_status": "aprovado",
    "order_payment_type": "CreditCard",
    "customer_id": 7,
    "customer_firstname": "Leandro",
    "customer_lastname": "Silva",
    "customer_email": "leandro@example.com",
    "customer_phone": "11999999999",
    "customer_site": {
      "id": 1470
    }
  }
}
```

---

## 4. Custom Content

**Detection:** `data.order_id` present, no `data.order_total_products`

Minimal payload — only the fields explicitly configured in the Appmax dashboard. The example below shows a common minimal setup.

**OrderApproved**
```json
{
  "event": "OrderApproved",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_status": "aprovado",
    "order_total": 267.48
  }
}
```

**OrderBilletCreated**
```json
{
  "event": "OrderBilletCreated",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_status": "pendente",
    "order_total": 267.48,
    "billet_url": "https://boleto.example.com/abc123"
  }
}
```

**OrderPixCreated**
```json
{
  "event": "OrderPixCreated",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_status": "pendente",
    "order_total": 267.48,
    "pix_code": "00020126330014br.gov.bcb.pix0111999999999952040000530398654041.005802BR5913Leandro Silva6008Sao Paulo62070503***6304ABCD"
  }
}
```

**OrderRefund**
```json
{
  "event": "OrderRefund",
  "event_type": "",
  "data": {
    "order_id": 12844,
    "order_status": "estornado",
    "order_total": 267.48
  }
}
```

---

## 5. Old Legacy

**Detection:** `event_type == "order"`

Pre-configurable-model format. Snake_case event names. `event_type` is always `"order"`. Still handled for backwards compatibility.

**order_approved**
```json
{
  "event": "order_approved",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_paid**
```json
{
  "event": "order_paid",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_billet_created**
```json
{
  "event": "order_billet_created",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_pix_created**
```json
{
  "event": "order_pix_created",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_paid_by_pix**
```json
{
  "event": "order_paid_by_pix",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_refund**
```json
{
  "event": "order_refund",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

**order_pix_expired**
```json
{
  "event": "order_pix_expired",
  "event_type": "order",
  "data": {
    "order_id": 12844
  }
}
```

---

## Event Reference

| Event | Model(s) | Category | Order ID field | Mapped status |
|-------|----------|----------|----------------|---------------|
| `OrderApproved` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `aprovado` |
| `OrderAuthorized` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `autorizado` |
| `OrderPaid` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `aprovado` |
| `OrderBilletCreated` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `pendente` |
| `OrderBilletOverdue` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `cancelado` |
| `OrderPixCreated` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `pendente` |
| `OrderPaidByPix` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `aprovado` |
| `OrderPixExpired` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `cancelado` |
| `OrderPendingIntegration` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `pendente_integracao` |
| `OrderIntegrated` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `integrado` |
| `OrderRefund` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `estornado` |
| `OrderChargeBackInTreatment` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `chargeback_em_tratativa` |
| `OrderUpSold` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | `aprovado` |
| `OrderPartialRefund` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | *(not mapped)* |
| `OrderChargeBackGain` | Standard, Standard with Meta, Two-Level Flat, Custom Content | Order | `data.id` / `data.order_id` | *(not mapped)* |
| `CreatedSubscription` | Standard, Standard with Meta, Two-Level Flat | Subscription | `data.id` / `data.order_id` | `aprovado` |
| `SubscriptionCancellationEvent` | Standard | Subscription | — (no-op) | — |
| `SubscriptionDelayedEvent` | Standard | Subscription | — (no-op) | — |
| `CustomerCreated` | Standard | Customer | — (no-op) | — |
| `CustomerInterested` | Standard | Customer | — (no-op) | — |
| `CustomerContacted` | Standard | Customer | — (no-op) | — |
| `order_authorized` | Old Legacy | Order | `data.order_id` | `autorizado` |
| `order_authorized_with_delay` | Old Legacy | Order | `data.order_id` | `autorizado` |
| `order_approved` | Old Legacy | Order | `data.order_id` | `aprovado` |
| `order_billet_created` | Old Legacy | Order | `data.order_id` | `pendente` |
| `order_paid` | Old Legacy | Order | `data.order_id` | `aprovado` |
| `order_pending_integration` | Old Legacy | Order | `data.order_id` | `pendente_integracao` |
| `order_refund` | Old Legacy | Order | `data.order_id` | `estornado` |
| `order_pix_created` | Old Legacy | Order | `data.order_id` | `pendente` |
| `order_paid_by_pix` | Old Legacy | Order | `data.order_id` | `aprovado` |
| `order_pix_expired` | Old Legacy | Order | `data.order_id` | `cancelado` |
| `order_integrated` | Old Legacy | Order | `data.order_id` | `integrado` |
| `order_billet_overdue` | Old Legacy | Order | `data.order_id` | `cancelado` |
| `order_chargeback_in_treatment` | Old Legacy | Order | `data.order_id` | `chargeback_em_tratativa` |
| `order_up_sold` | Old Legacy | Order | `data.order_id` | `aprovado` |
| `payment_not_authorized` | Old Legacy | Order | `data.order_id` | `cancelado` |
| `payment_authorized_with_delay` | Old Legacy | Order | `data.order_id` | `autorizado` |
| `split_orders` | Old Legacy | Order | `data.order_id` | `aprovado` |
| `customer_created` | Old Legacy | Customer | — (no-op) | — |
| `customer_interested` | Old Legacy | Customer | — (no-op) | — |
| `customer_contacted` | Old Legacy | Customer | — (no-op) | — |
| `subscription_cancelation` | Old Legacy | Subscription | — (no-op) | — |
| `subscription_delayed` | Old Legacy | Subscription | — (no-op) | — |
