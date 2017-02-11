# MP3 Tag Fixer

This command line utility provides a flexible way to set the title meta data of a MP3 file based on the file's name. 

# Why?

I use Downcast, which is a great podcast player, to listen to audio books. Downcast treats each mp3 file in an audio book as a podcast and provides options for ordering them by title. Many mp3 audio books do not have title metadata that sorts nicely. Often the title is the same for each file. 

# How

The idea is to extract the title and index number for each file from the file name, which almost always has the index number and title in it in some fashion. Regular expressions with sub-expression matching are employed. Once properly extracted they are put into a template and that is set as the mp3 title. 

**Example file name:**

`Author Name - Book Title Ch34.mp3`

**Index pattern**

`([0-9]+)`

**Name pattern**

`- (A-Za-z0-9 ]+) Ch.*`

**Title template**

`{{ .Index }} - {{ .Name }}`

**RESULT**

`34 - Book Title`

The sub expressions (everything between the parenthesis) are what tagfix will use. It must match exactly once or the file is skipped. 

# Command Help

```Usage of tagfix:
  -idxpat string
        regex pattern to extract file index string (default "([0-9]+)")
  -mp3dir string
        directory of mp3 files
  -namepat string
        regex pattern to extract name of item (default "- ([A-Za-z0-9 ]+) Ch.*")
  -quiet
        don't make a lot of noise
  -titletemp string
        a golang text template for the id3 title to be set to (default "{{ .Index }} - {{ .Name }}")
  -trial
        try things out without modifying the mp3 file```
        
# Warnings

The mikkyang-id3 library seems to get rid of any images stored in the metadata and that is not important for my use cases.

