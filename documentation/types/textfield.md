# Textfield

A textfield is an element, well, that you type text in. In the lua api, there is not much you can do with it besides doing x.text to get/set the value of the text in it.

ex.
```jtl
>id="mytextfield">textfield>I have not changed yet;
```

Styles:
    width: Stretches the element on the width in px. EX: `width: 50;`
    height: Stretches the element on the height in px. EX: `height: 50;`
    color: Sets the color of the element in rgba (red, green, blue, alpha). EX: `color: 0, 0, 0, 255`
    font-family: Selects a font to use for displaying. Read avalible fonts in documentation/fonts.txt. EX: `font-family: JetBrainsMono;`

Lua Attributes:
    .text, string: gets/sets text of the object. Ex:
```lua
local x = document.get("#mytextfield")
print(x.text)
x.text = "i have magically changed"
```

elem.class: String of class name.
elem.identifier: String of id.