# Login by SSH

A demonstration of using SSH to register and authenticate users.

Run the server and try:

```sh
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p2222 -l<user> localhost
```

where `<user>` is your desired username. Do not omit `StrictHostKeyChecking` and `UserKnownHostsFile` options on your local machine to prevent putting junk in your SSH config files.
