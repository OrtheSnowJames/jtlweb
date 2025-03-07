# Button
A simple button that people know how to click probably. Use this for forms, games or whatever. Interacts with the lua api via on event and text.

ex.
```jtl
>class="buttonclass">button>Click me!;
```

Styles:
    width: Stretches the element on the width in px. EX: `width: 50;`
    height: Stretches the element on the height in px. EX: `height: 50;`
    color: Sets the color of the element in rgba (red, green, blue, alpha). EX: `color: 0, 0, 0, 255`
    border-color: Sets the border color of the element in rgba (red, green, blue, alpha). EX: `border-color: 0, 0, 0, 255`
    font-family: Selects a font to use for displaying. Read avalible fonts in documentation/fonts.txt. EX: `font-family: JetBrainsMono;`
    margin: Makes space around the element in px to leave blank space. EX: `margin: 50;`
    padding: A fancy name for margin (yes, it is the same thing). EX: `padding: 50;`

Lua Attributes:
    event.click: Happens on a click. Ex:
```lua
function handler()
    print("Button Clicked!")
end 

document.onEvent(".buttonclass", "click", [[handler()]])
```
event.clickrepeat: Happens on every frame of the button being clicked. Ex:

```lua
function handleralso()
    print("spam")
end 

document.onEvent(".buttonclass", "clickrepeat", [[handleralso()]])
```

elem.class: String of class name.
elem.identifier: String of id.