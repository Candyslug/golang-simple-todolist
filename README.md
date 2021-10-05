# golang-simple-todolist
prototype todolist
<p>Made in golang, using the standard library</p>

<br>

<p>Simple todolist api, made quickly to get used to the language</p>

<br>

<ul>
  <li>"/items" (GET) will show all items as json</p></li>
  <li>"/items/x" (GET) will show item at x (int) index</p></li>
  <br>
  <li>"/items/add" (POST json struct) will add item to list</p></li>
  <li>"/items/del" (POST index struct) will remove item index x from list</p></li>
  <i>"/mark" (POST indextoggle struct) will toggle completion on item index x on list</p></li>
</ul>
