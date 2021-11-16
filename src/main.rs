use log::{info, warn};

fn main() {
    log::log_enabled!(log::Level::Debug);
    println!("Hello, world!");
    info!("info...");
}
