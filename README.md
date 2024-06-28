<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h3 align="center">Vizzy</h3>

  <p align="center">
    Web application written in Go, HTMX, and Tailwind to assist with the sharing of nmap scans, credentials, and other information during pentests or red team engagements.
    <br />
    <br />
    <a href="#installation">Install</a>
    ·
    <a href="https://github.com/ziggoon/vizzy/issues">Report Bug</a>
    ·
    <a href="https://github.com/ziggoon/vizzy/issues">Request Feature</a>
  </p>
</div>


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#built-with">Built With</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



### Built With
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Instructions on how to run Vizzy locally or in a container using the provided Dockerfile

### Prerequisites

The only required software is Golang, which can be downloaded at the following url
  ```sh
  https://go.dev/dl/
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/ziggoon/vizzy
   ```
2. Change directory
   ```sh
   cd vizzy
   ```
3. Build executable
   ```sh
   go build
   ```
4. Run executable
   ```sh
   ./vizzy
   ```

---
_These steps are for Docker users only_

1. Clone the repo
   ```sh
   git clone https://github.com/ziggoon/vizzy
   ```
2. Change directory
   ```sh
   cd vizzy
   ```
3. Build docker container
   ```sh
   docker build -t vizzy .
   ```

4. Run docker container
    ```sh
    docker run -p 42069:42069 --name vizzy vizzy
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

- [x] Add authentication
- [ ] Add admin features (create users, delete users, etc..)
- [ ] tbd

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the GPLv3 License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

ziggoon - [twitter](https://twitter.com/0xzig)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
