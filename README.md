# go-components

golang components

## list

| Name      | package                                    | function                                                         | distributed system support | dependency |
|-----------|--------------------------------------------|------------------------------------------------------------------|----------------------------|------------|
| Audit     | github.com/junqirao/gocomponents/audit     | audit log                                                        | √                          | -          |
| Grace     | github.com/junqirao/gocomponents/grace     | graceful exit                                                    | √                          | -          |
| KVDB      | github.com/junqirao/gocomponents/kvdb      | key-value based database, derivative usage: message bus, storage | √                          | -          |
| Launcher  | github.com/junqirao/gocomponents/launcher  | stander launch helper                                            | √                          | -          |
| MFA       | github.com/junqirao/gocomponents/mfa       | MFA utils                                                        | √                          | -          |
| Procedure | github.com/junqirao/gocomponents/procedure | input/output abstract and configurable                           | √                          | -          |
| Registry  | github.com/junqirao/gocomponents/registry  | registry for service discovery                                   | √                          | KVDB       |
| Response  | github.com/junqirao/gocomponents/response  | goframe response enhancement                                     | √                          | -          |
| Security  | github.com/junqirao/gocomponents/security  | rsa security provider                                            | √                          | -          |
| Storage   | github.com/junqirao/gocomponents/storage   | object storage                                                   | √                          | -          |
| Task      | github.com/junqirao/gocomponents/task      | async task helper                                                | √                          | KVDB       |
| Updater   | github.com/junqirao/gocomponents/updater   | service update tasks, update something on start                  | √                          | -          |
