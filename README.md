# PAYCHECK

A small script to download all your paychecks from cloudserver provider.

## CONFIG
Check if you have the following file: `~/.config/paycheck/config.json`
otherwise create it and copy inside `config.json` the following code:
```json
{
  "svcloud_username": "your-username",
  "svcloud_password": "your-password",
  "svcloud_base_url": "https://my-base-url",
  "svcloud_endpoint": "/my-endpoint",
  "box_token"       : "The BOX auth token",
  "box_target_id"   : "The folder id of box"
}
```
NOTE:
Make sure to set the correct permissions for config.json. You
should consider to set `600` the permissions of this file

### Box configuration
This script use box.com service to store the documents downloaded.
To start using BOX, you must first of all signup to the service and then create an APP with
authentication method: 
```
App Token (Server Authentication)
```
[https://developer.box.com/guides/authentication/app-token/](More detail here)

Then you must create a subfolder inside the main root and share it with the service account created
during app creation process. (you can find this info in: *general settings > Service account info*)
Under the section **Service account info** you can find an email. Share the folder adding only upload 
privilege with this user (fake).

### GPG configuration
If you want to use the GGP encryption you must first of all create a new key if you
doesn't have one. To create a new key run the following commands:
```bash
# Create new key following the the TUI wizard
gpg --full-generate-key

# After the creation process export the public key armored version
gpg --export --armor [KEY ID] > /tmp/paycheck.asc

# Put the file just created inside the conf folder:
# ~/.config/paycheck
mv /tmp/paycheck.asc ~/.config/paycheck

# Now you are able to encrypt the files to upload
```

## USAGE
After finish to fill the `config.json`, you can simply run the following command:
```bash
paycheck
```

### Usage options
The following flags are available:
- `-c` with this flag you can apply the gpg encryption.
- `-u` with this flag you can also upload the file encrypted to your **box account**.


## RUN
If you have `GO` installed, download the repo, set teh config and then run `go run main.go`.
Otherwise download the latest binaries and run with `./paycheck`.
You can also save the binary inside a bin folder to execute it globally.

