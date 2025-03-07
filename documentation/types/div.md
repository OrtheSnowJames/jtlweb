# Div

A div is a container element that can hold other elements. It inherits styles to its children.

Ex.
```jtl
>>>DOCTYPE=JTL

>>>BEGIN;
    >style="margin-left=50">div>
        >class="textclass" noattribute="true">p>helo;
        >class="textclass" noattribute="true">p>helo;
>>>END;
```

Lua attributes:
elem.chilren: Returns children in array lua table form.
