# i3 workspace auto rename

Automatically display the opened application names in the workspaces (using the i3 rename workspace feature).

![screenshot](screenshot.png)

## i3 config

It is important to define workspace with numbers and no initial name in the i3 config file to keep the key binding to switch workspace.

```bash
set $ws1 1
set $ws2 2
set $ws3 3
...

workspace $ws1 output DP-1-2
workspace $ws2 output DP-1-2
workspace $ws3 output DP-1-2
...

bindsym $mod+ampersand workspace number 1
bindsym $mod+eacute workspace number 2
bindsym $mod+quotedbl workspace number 3
```

## installation

Drop the binary somewhere and add to you i3 config file:

```bash
exec /path/to/my/goi3autowname
```

## run options

`--debug`: debug mode
`--mapf`: full path (with file name) to the `goi3autowname.json` file containing a `class`<>`name` mapping.