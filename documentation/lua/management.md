# add, create, insert, remove and replace
Add:
Inserts an element at the top of the document.

Create:
Creates a temporary element for editing and adding.

Insert:
Inserts the document at a specific position with the lua offset (lua offset = first element of vec will be at vec[ 1 ])

Remove:
Removes an item with a query selector like when you get an item

Replace:
Replaces an item with a different item based on a query selector

EX:

```lua
local x = document.create("p")
x.text = "Hello, world!"
x.class = "coolclass"
document.add(x)

local y = document.create("p")
y.text = "Hello, World!"
y.class = "evencoolerclass"
y.identifier = "yid"
document.insert(y, 1) -- inserts at the top of the document

local z = document.create("p")
z.text = "hello"
z.class = "evenmorecoolerclass"
z.identifier = "zid"
document.replace(".coolclass", z) -- replaces x with z

document.remove(".evenmorecoolerclass") -- removes z
document.remove("#yid") -- removes y
```