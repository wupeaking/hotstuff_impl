mod p2p;
pub mod config;
use config::HotStuffConf;
use anyhow::Error;


pub trait HostStuffModule {
    type service;
    fn new(cfg :&HotStuffConf) -> Self::service;

    // 启动
    fn start(&mut self) -> Result<(), Error>;

    // 停止
    fn stop(&mut self) -> Result<(), Error>;
}