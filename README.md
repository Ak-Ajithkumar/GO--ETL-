# GO--ETL-

# Image Processing ETL

This Go program processes images by converting them to grayscale and creates a ZIP archive of the processed images. It uses concurrency to efficiently handle multiple files.

## Features

- Convert images to grayscale.
- Process images concurrently.
- Create ZIP archives of processed images.
- Handle unique filenames to avoid overwriting.

## Requirements

- Go (version latest)

## Installation

1. **Clone the repository:**


2. **Install the Dependency:**
     
- go get github.com/disintegration/imaging
- go get github.com/mholt/archiver/v3

3. **Make input file :**

- Add images to the input directory.

4. **Run Command:**

- go run main.go