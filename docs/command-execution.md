# Command Execution

> **Note:** The hook runner and interactive shell are disabled by default. Only enable them if you understand the security implications.


## Hook Runner

The hook runner is a feature that enables you to execute any shell command you want before or after a certain event. Right now, these are the events:

* Copy
* Rename
* Upload
* Delete
* Save

Also, during the execution of the commands set for those hooks, there will be some environment variables available to help you perform your commands:

* `FILE` with the full absolute path to the changed file.
* `SCOPE` with the path to user's scope.
* `TRIGGER` with the name of the event.
* `USERNAME` with the user's username.
* `DESTINATION` with the absolute path to the destination. Only used for **copy** and **rename.**

At this moment, you can edit the commands via the command line interface, using the following commands \(please check the flag `--help` to know more about them\):

```bash
filebrowser cmds add before_copy "echo $FILE"
filebrowser cmds rm before_copy 0
filebrowser cmds ls
```

Or you can use the web interface to manage them via **Settings** → **Global Settings**.

### How the variables are filled in

* **Without a shell** (the default), File Browser replaces `$FILE`, `${FILE}` and the other variables inside each argument before running the program. No shell is involved, so a file name is always passed as a single argument.
* **With a shell** (the `shell` setting, for example `sh -c`), the hook is handed to the shell as a script and the shell expands the variables itself from the environment. Quote them so names with spaces stay together:

  ```bash
  filebrowser cmds add after_upload 'echo "$FILE" >> /var/log/uploads.log'
  ```

  Use your shell's own syntax: `"$FILE"` for `sh`/`bash`, `%FILE%` for `cmd`, `$env:FILE` for PowerShell.

> **Note:** File names are never pasted into a shell script. Earlier versions did, so an uploaded file named like `a;rm -rf ~` could run commands when a shell was configured. If you wrote hooks for `cmd` or PowerShell using `$FILE`, switch them to that shell's syntax.

## Interactive Shell

Within File Browser you can toggle the shell (`< >` icon at the top right) and this will open a shell command window at the bottom of the screen. This functionality can be turned on using the environment variable `FB_DISABLE_EXEC=false` or the flag `--disable-exec=false`.

By default no commands are available as the command list is empty. To enable commands these need to either be done on a per-user basis (including for the Admin user).

You can do this by adding them in Settings > User Management > (edit user) > Commands or to *apply to all new users created from that point forward* they can be set in Settings > Global Settings

> **Note:** The allowed-commands list only checks the **first word** of what the user types. When a `shell` is configured, the whole line goes to that shell, so a user allowed to run `ls` can also run `ls; anything`. Leave `shell` empty if you rely on the list to restrict users, or only give the terminal to people you would trust with a real shell.

Closing the terminal tab (or leaving the page) stops the running command.

> **Note:** If using a proxy manager then remember to enable websockets support for the File Browser proxy

> **Note:** If using Docker and you want to add a new command that is not in the base image then you will need to build a custom Docker image using `filebrowser` as a base image.  For example to add 7z:
>
> ```docker
> FROM filebrowser
> RUN sudo apt install p7zip-full
> ```
