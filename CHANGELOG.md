# Changelog

All notable changes to this crate will be documented in this file.

The format is based on [Keep a Changelog], and this crate adheres to [Semantic
Versioning].

## 0.1.1 - 2021-08-17

* Lazily connect to the PostgreSQL database the first time a resource managed
  by the provider is created or destroyed. This makes it possible to run
  `pulumi preview` when the database connection is not available.

## 0.1.0 - 2021-07-06

Initial release.

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
