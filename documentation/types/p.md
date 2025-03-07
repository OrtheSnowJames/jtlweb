# P

The p element is an element holding paragraph text, although it can be used for other things.

Ex.
```jtl
>noattribute="true">p>Hello, World!;
```

This defines a simple "p" element that shows the text Hello, World! upon the view.

Styles:
    font-family: Selects a font to use for displaying. Read avalible fonts in documentation/fonts.txt. EX: `font-family: JetBrainsMono;`
    width: Stretches the element on the width in px. EX: `width: 50;`
    height: Stretches the element on the height in px. EX: `height: 50;`
    margin: Makes space around the element in px to leave blank space. EX: `margin: 50;`
    margin-left: Sets the left margin of the element in px. EX: `margin-left: 10;`
    margin-right: Sets the right margin of the element in px. EX: `margin-right: 10;`
    margin-up: Sets the top margin of the element in px. EX: `margin-up: 10;`
    margin-down: Sets the bottom margin of the element in px. EX: `margin-down: 10;`
    center: This is not recommended, but it puts the page in the center of the screen size. Doesn't scroll with you. EX: `center: true;`

Lua Attributes:
    .text, string: gets/sets text of the object. Ex:
```lua
local x = document.get("p") -- assuming you want to get the first p element in the document
print(x.text)
x.text = "i have magically changed"
```

elem.class: String of class name.
elem.identifier: String of id.