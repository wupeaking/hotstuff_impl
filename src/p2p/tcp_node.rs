use crate::HostStuffModule;
use crate::config::{HotStuffConf, P2PConf};
use tokio::net::{TcpListener, TcpStream};
use std::collections::HashMap;


// tcp协议的P2P节点实现
struct TCPNode<'a> {
    cfg : &'a P2PConf,
    peers: HashMap<String,TcpStream>
}

impl<'a> HostStuffModule for TCPNode<'a> {
    type Service = Self;
    fn new(cfg :&HotStuffConf) -> Self {
        return TCPNode{
            cfg: &cfg.p2p_cfg,
            peers: HashMap::new(),
        };
    }

    // todo::
    fn start(&mut self) -> Result<(), anyhow::Error> {
        Ok(())
    }

    // todo::
    fn stop(&mut self) -> Result<(), anyhow::Error> {
        Ok(())
    }

}
