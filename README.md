<p align="center">
    <h3 align="center">Nessus on Demand</h3>
    <p align="center">
        A just-in-time Nessus solution for modern penetration testing teams.
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

Nessus on Demand (NoD) is a handy CLI utility for managing Nessus instances in AWS. Built using Terraform, NoD safely bootstraps scanning infrastructure on a per-need basis.

## Installation and Usage

To install Nessus on Demand, use the `go install` command.

```bash
$ go install github.com/prdngr/nessus-on-demand@latest
```

Once installed, the easiest way of spinning up an instance using NoD looks as follows:

> [!IMPORTANT]
> The command shown below will use your AWS CLI default profile. If you want to use another profile instead, set the [`AWS_PROFILE`](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-using-profiles) environment variable accordingly.

```bash
$ nessus-on-demand deployment create --auto-ip --region eu-central-1

    ███╗   ██╗ ██████╗ ██████╗
    ████╗  ██║██╔═══██╗██╔══██╗
    ██╔██╗ ██║██║   ██║██║  ██║
    ██║╚██╗██║██║   ██║██║  ██║
    ██║ ╚████║╚██████╔╝██████╔╝
    ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝

NoD initialized ✅
Nessus deployed ✅
Deployment details gathered ✅

Deployment Summary
------------------

Deployment ID: 90d6c9be-edf0-4ee9-ac0b-c7ed701c6f70
Nessus Interface: https://3.176.102.79:8834
Allowed IP Address: 118.146.33.52

Next Steps
----------

▶ Open the Nessus interface in your browser, sign up, and activate your license.
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

4. Open a [pull request](https://github.com/prdngr/nessus-on-demand/pulls).

## License

Distributed under the MIT License. See `LICENSE` for more information.
