# A binary to download Wordpress for the given version and language

Just clone the repo, cd into it and build the code:

```shell
./build.sh
```

Then run the generated executable using the **country code** you want. By default is _en_:

```shell
./get-lastest-wp
```

**country code** is implicitly set to _en_.

```shell
./get-lastest-wp -country es
```

**country code** is set to _es_.

If the **language code** for the given **country code** is not found, you can manually indicate both using the _-lang_ flag also. For example:

```shell
./get-lastest-wp -country zz -lang zz_ZZ
```

By default the latest version is downloaded, though there's also another flag to indicate a **version** other than the latest one:

```shell
./get-lastest-wp -version 4.9.5
```
