return errors.Wrap(err, "request failed")
 
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

https://www.reddit.com/r/golang/comments/86mpvg/how_do_you_deal_with_rdbms_relationships_in_go/
https://tristangio.fr/2017/03/go-simplify-sql-with-sqlx/
https://jmoiron.github.io/sqlx/

http://www.sqlservertutorial.net/sql-server-basics/sql-server-left-join/

// Dont use interface
https://www.calhoun.io/interfaces-define-requirements/

https://github.com/katzien/go-structure-examples/tree/master/domain-hex
https://github.com/benbjohnson/wtf/tree/http
https://programmingwithmosh.com/net/common-mistakes-with-the-repository-pattern/

// Put order and items in same repo, use tx have method per op

defer tx.Rollback()

public interface BusinessRuleGateway {
  Something getSomething(String id);
  void startTransaction();
  void saveSomething(Something thing);
  void endTransaction();
}

//open transaction and set in participants
tx := openTransaction()
ur := NewUserRepository(tx)
ir := NewImageRepository(tx)
//keep user and image datas
err0 := ur.Keep(userData)
err1 := ir.Keep(imageData)
//decision
if err0 != nil || err1 != nil {
  tx.Rollback()
  return
}
tx.Commit()

SELECT TOP(5) ProductID, SUM(Quantity) AS TotalQuantity
FROM order_items
GROUP BY ProductID
ORDER BY SUM(Quantity) DESC;

// db options with wrap
https://hackernoon.com/how-to-work-with-databases-in-golang-33b002aa8c47

// mongo paginate
https://github.com/saeedghx68/fast-relay-pagination/blob/6eb6af05ae96df174fd18b220e7904ef69648e9a/src/index.js