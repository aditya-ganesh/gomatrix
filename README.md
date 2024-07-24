# gomatrix

Yet another matrix rain implementation, this time in go (It's not even the [first one](https://github.com/GeertJohan/gomatrix), which has been used as a reference).

##

Command line args

| **Flag** | **Type** | **Default** | **Description**                                                                                                                                                          |
|----------|----------|-------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| c        | string   | None        | Raindrop colour. If supplied, must be a colour name that [Tcell supports](https://github.com/gdamore/tcell/blob/88b9c25c3c5ee48b611dfeca9a2e9cf07812c35e/color.go#L851)  |
| min      | float64  | 0.2         | Minimum raindrop size given as a fraction of the screen height                                                                                                           |
| max      | float64  | 0.5         | Maximum raindrop size given as a fraction of the screen height                                                                                                           |
| r        | float64  | 0.1         | Refresh interval in seconds                                                                                                                                              |