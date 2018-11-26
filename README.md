# i3 workspace auto rename

## i3 confiu

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

## planned

Add a appclass<->name config file to display custom workspace names.