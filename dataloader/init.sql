-- Adminer 4.7.1 PostgreSQL dump

\connect "default";

DROP TABLE IF EXISTS "items";
DROP SEQUENCE IF EXISTS items_id_seq;
CREATE SEQUENCE items_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."items" (
                                  "id" integer DEFAULT nextval('items_id_seq') NOT NULL,
                                  "title" text,
                                  CONSTRAINT "items_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "items" ("id", "title") VALUES
(1,	'item1'),
(2,	'item2'),
(3,	'item3');

DROP TABLE IF EXISTS "items_tags";
CREATE TABLE "public"."items_tags" (
                                       "item_id" integer NOT NULL,
                                       "tag_id" integer NOT NULL,
                                       CONSTRAINT "items_tags_pkey" PRIMARY KEY ("item_id", "tag_id"),
                                       CONSTRAINT "items_tags_item_id_fkey" FOREIGN KEY (item_id) REFERENCES items(id) NOT DEFERRABLE,
                                       CONSTRAINT "items_tags_tag_id_fkey" FOREIGN KEY (tag_id) REFERENCES tags(id) NOT DEFERRABLE
) WITH (oids = false);

INSERT INTO "items_tags" ("item_id", "tag_id") VALUES
(1,	1),
(1,	2),
(1,	3),
(1,	4),
(2,	1),
(2,	2),
(2,	3),
(3,	4),
(3,	5);

DROP TABLE IF EXISTS "order_items";
CREATE TABLE "public"."order_items" (
                                        "orderid" character varying(36) NOT NULL,
                                        "productid" character varying(36) NOT NULL,
                                        "name" character varying(50) NOT NULL,
                                        "uom" character varying(3) NOT NULL,
                                        "requested" integer NOT NULL,
                                        "received" integer NOT NULL,
                                        "remaining" integer NOT NULL,
                                        "status" integer NOT NULL,
                                        "ponumber" character varying(20) NOT NULL,
                                        "dateadded" integer NOT NULL,
                                        CONSTRAINT "order_items_orderid_productid_unique" UNIQUE ("orderid", "productid"),
                                        CONSTRAINT "order_items_orderid_fkey" FOREIGN KEY (orderid) REFERENCES orders(orderid) ON DELETE CASCADE NOT DEFERRABLE
) WITH (oids = false);

CREATE INDEX "order_items_orderid_index" ON "public"."order_items" USING btree ("orderid");


DROP TABLE IF EXISTS "orders";
CREATE TABLE "public"."orders" (
                                   "orderid" character varying(36) NOT NULL,
                                   "projectid" character varying(36) NOT NULL,
                                   "sentdate" integer NOT NULL,
                                   "status" integer NOT NULL,
                                   "project_name" character varying(80) NOT NULL,
                                   "foreman_email" character varying(80) NOT NULL,
                                   "comments" character varying(3000) NOT NULL,
                                   CONSTRAINT "orders_pk" PRIMARY KEY ("orderid"),
                                   CONSTRAINT "orders_projectid_fkey" FOREIGN KEY (projectid) REFERENCES projects(id) NOT DEFERRABLE
) WITH (oids = false);

INSERT INTO "orders" ("orderid", "projectid", "sentdate", "status", "project_name", "foreman_email", "comments") VALUES
('170',	'projectID01',	170,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('160',	'projectID01',	160,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('150',	'projectID01',	150,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('140',	'projectID01',	140,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('130',	'projectID02',	130,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('120',	'projectID01',	120,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('110',	'projectID01',	110,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('100',	'projectID02',	100,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('90',	'projectID01',	90,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('80',	'projectID01',	80,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('70',	'projectID02',	70,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('50',	'projectID01',	50,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('40',	'projectID01',	40,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('30',	'projectID01',	30,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('20',	'projectID01',	20,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!'),
('180',	'projectID03',	180,	1,	'Argyle Secondary School',	'mwright@plan-group.com',	'Nothing today, thanks!');

DROP TABLE IF EXISTS "person";
CREATE TABLE "public"."person" (
                                   "first_name" text,
                                   "last_name" text,
                                   "email" text
) WITH (oids = false);

INSERT INTO "person" ("first_name", "last_name", "email") VALUES
('Jason',	'Moiron',	'jmoiron@jmoiron.net'),
('John',	'Doe',	'johndoeDNE@gmail.net'),
('Jane',	'Citizen',	'jane.citzen@example.com');

DROP TABLE IF EXISTS "place";
CREATE TABLE "public"."place" (
                                  "country" text,
                                  "city" text,
                                  "telcode" integer
) WITH (oids = false);

INSERT INTO "place" ("country", "city", "telcode") VALUES
('United States',	'New York',	1),
('Hong Kong',	NULL,	852),
('Singapore',	NULL,	65);

DROP TABLE IF EXISTS "projects";
CREATE TABLE "public"."projects" (
                                     "id" text NOT NULL,
                                     "name" text NOT NULL,
                                     CONSTRAINT "projects_id" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "projects" ("id", "name") VALUES
('projectID01',	'Project 01'),
('projectID02',	'Project 02'),
('projectID03',	'Project 03');

DROP TABLE IF EXISTS "tags";
DROP SEQUENCE IF EXISTS tags_id_seq;
CREATE SEQUENCE tags_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."tags" (
                                 "id" integer DEFAULT nextval('tags_id_seq') NOT NULL,
                                 "title" text,
                                 CONSTRAINT "tags_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "tags" ("id", "title") VALUES
(1,	'tag1'),
(2,	'tag2'),
(3,	'tag3'),
(4,	'tag4'),
(5,	'tag5');

-- 2019-05-20 20:07:01.167072+00
