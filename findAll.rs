use std::{env, fs, io::{self, BufRead}, path::Path};

fn search_in_file(file_path: &Path, pattern: &str) {
    if let Ok(file) = fs::File::open(file_path) {
        let reader = io::BufReader::new(file);
        for (line_num, line) in reader.lines().enumerate() {
            if let Ok(content) = line {
                if content.contains(pattern) {
                    let highlighted = content.replace(pattern, &format!("\x1b[31m{}\x1b[0m", pattern));
                    println!("{}:{}: {}", file_path.display(), line_num + 1, highlighted);
                }
            }
        }
    }
}

fn search_in_directory(dir: &Path, pattern: &str) {
    if let Ok(entries) = fs::read_dir(dir) {
        for entry in entries.flatten() {
            let path = entry.path();
            if path.is_file() {
                search_in_file(&path, pattern);
            } else if path.is_dir() {
                search_in_directory(&path, pattern);
            }
        }
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 3 {
        eprintln!("Usage: findAll <directory> <pattern>");
        return;
    }

    let dir = Path::new(&args[1]);
    let pattern = &args[2];
    
    if dir.is_dir() {
        search_in_directory(dir, pattern);
    } else {
        eprintln!("Error: '{}' is not a directory", dir.display());
    }
}

