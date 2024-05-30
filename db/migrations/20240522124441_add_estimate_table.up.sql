
CREATE TABLE IF NOT EXISTS "estimate" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "userId" uuid NOT NULL,
    "userLat" REAL NOT NULL,
    "userLon" REAL NOT NULL,
    "totalPrice" int NOT NULL,
    -- time in second, then convert to minutes in code
    "estimateDeliveryTime" int NOT NULL,
    -- if null, not checkout yet, NOT FK
    "orderId" uuid UNIQUE DEFAULT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW()),
    "updatedAt" timestamp NOT NULL DEFAULT(NOW())
);
 
CREATE TABLE IF NOT EXISTS "estimateOrder"(
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "estimateId" uuid NOT NULL,
    "isStarting" bool DEFAULT FALSE,
    "merchantId" uuid NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW())
);

CREATE TABLE IF NOT EXISTS "estimateOrderItem" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "estimateOrderId" uuid NOT NULL,
    "itemId" uuid NOT NULL,
    "quantity" int NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW())
);

ALTER TABLE "estimate" ADD FOREIGN KEY ("userId") REFERENCES "user" ("id");
ALTER TABLE "estimateOrder" ADD FOREIGN KEY ("estimateId") REFERENCES "estimate" ("id");
ALTER TABLE "estimateOrder" ADD FOREIGN KEY ("merchantId") REFERENCES "merchant" ("id");
ALTER TABLE "estimateOrderItem" ADD FOREIGN KEY ("estimateOrderId") REFERENCES "estimateOrder" ("id");
ALTER TABLE "estimateOrderItem" ADD FOREIGN KEY ("itemId") REFERENCES "item" ("id");