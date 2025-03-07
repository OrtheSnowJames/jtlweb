Some other useful call apis

## onFrame
(Not Recommended because lua is kind of slow) Registers a function to be executed every frame.
Ex:
```lua
function handler()
    print("hi")
end
document.onFrame([[handler()]])
```

## Get
Uses a tag to get an object from the page.
Ex:
```lua
local x = document.get("#par") -- gets element with id of par
local y = document.get(".par") -- gets element with class of par
local z = document.get("p") -- gets first element with type of p

print(x.text + y.text + z.text)
```

## requestFrame
Registers a function to be executed on the next frame only.
Ex:
```lua
-- this is how to make a proper loop
function loop()
    document.requestFrame([[loop()]])
end

loop()
```

## objects
Gets all objects from the page.
Ex:
```lua
local objs = document.objects()
for i, obj in ipairs(objs) do
    print(obj)
end
```

## update
Updates an element in the document.
Ex:
```lua
local elem = document.get("#par")
elem.text = "Updated text"
document.update(elem)
```

## whereis
Finds the position of an element with the lua offset.
Ex:
```lua
local elem = document.get("#par")
local position = document.whereis(elem)
print("Element position: " .. position)
```

## addStyle
Adds a style to an element.
Ex:
```lua
local elem = document.get("#par")
document.addStyle(elem, "color", "red")
```

## removeAllStyle
Removes all styles from an element.
Ex:
```lua
local elem = document.get("#par")
document.removeAllStyle(elem)
```

## fetch
Fetches data from a JTLTP URL.
Ex:
```lua
local tableOfResponse = document.fetch("jtltp://example.com", [[onResponse()]])

local responseStatus = tableOfResponse["JTLTP-STATUS"]
local docType = tableOfResponse["JTLTP-TYPE"]
local responseContent = tableOfResponse["JTLTP"]
```