# vcn - Testing vcn on windows

In order to stress test vcn on windows we used [machma](https://github.com/fd0/machma)

```bash
   dir /s /b *.* | machma.exe -p 4 -- vcn.exe --lc-api-key=my-cnil-api-key n {}
```
