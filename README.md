<p align="center">
    <h3 align="center">RedSky</h3>
    <p align="center">
        ☁️ The Calm Before the Breach
    </p>
</p>

<details>
<summary>Table of Contents</summary>

- [About The Project](#about-the-project)
- [Installation and Usage](#installation-and-usage)
- [Contributing](#contributing)
- [License](#license)

</details>

## About The Project

RedSky is a handy CLI utility for managing just-in-time offensive infrastructure in AWS. Currently, RedSky supports the following deployment types:

- Tenable Nessus (BYOL)
- Kali Linux standalone
- Kali Linux with Mythic C2

## Installation and Usage

To install RedSky, use the `go install` command.

```bash
$ go install github.com/prdngr/red-sky@latest
```

Once installed, the easiest way of spinning up a Kali instance using RedSky looks as follows:

```bash
$ red-sky create --type kali --region eu-west-1 --auto-ip

        ____           _______ __
       / __ \___  ____/ / ___// /____  __
      / /_/ / _ \/ __  /\__ \/ //_/ / / /
     / _, _/  __/ /_/ /___/ / ,< / /_/ /
    /_/ |_|\___/\__,_//____/_/|_|\__, /
                                /____/

✅ AWS session initialized

AWS Account Information
-----------------------

▶ Account: 444444444444 (red-sky)
▶ Caller ARN: arn:aws:iam::444444444444:user/red-sky

Press Enter to continue...

✅ RedSky initialized
✅ Deployment executed
✅ Deployment details gathered

Deployment Summary
------------------

▶ Deployment ID: d1505235-2e81-49b0-8bb6-3b1b76616b00
▶ Allowed IP Address: 42.42.42.42

Connection Details
------------------

▶ Use the following command to SSH into the Kali instance:
  $ ssh -i 'd1505235-2e81-49b0-8bb6-3b1b76616b00.pem' kali@62.62.62.62
```

## Contributing

The project is developed according to the [GitFlow workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) and it is encouraged to follow these [Git commit message guidelines](https://gist.github.com/robertpainsi/b632364184e70900af4ab688decf6f53).

1. Create your feature branch:

    ```console
    git checkout -b feature/<feature-name>
    ```

2. Commit your changes:

    ```console
    git commit -m '<commit-message>'
    ```

3. Push to the feature branch:

    ```console
    git push origin feature/<feature-name>
    ```

4. Open a [pull request](https://github.com/prdngr/red-sky/pulls).

## License

Distributed under the GNU GPLv3 License. See `LICENSE` for more information.
