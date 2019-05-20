realize start

Scrollbar on hover 

div {
  height: 100px;
  width: 50%;
  margin: 0 auto;
  overflow: hidden;
}

div:hover {
  overflow-y: scroll;
}

SELECT i.id, i.title, array_agg(i.title)
FROM items i
INNER JOIN items_tags it
ON it.item_id = i.id
INNER JOIN tags t
ON t.id = it.tag_id
GROUP BY i.id, i.title,

https://lorenstewart.me/2017/12/03/postgresqls-array_agg-function/
https://www.opsdash.com/blog/postgres-arrays-golang.html
https://stackoverflow.com/questions/44379851/get-postgresql-array-into-struct-with-structscan
https://github.com/jmoiron/sqlx/issues/168
s := strings.Split("a,b,c", ",")