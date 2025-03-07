local divelem = document.get("div")

local divelems = div.children -- returns an array table of all children of the div element

for i, elem in ipairs(divelems) do
    print(elem.text)
end