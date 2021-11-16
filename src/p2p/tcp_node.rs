use crate::HostStuffModule;
use crate::config::{HotStuffConf, P2PConf};
use tokio::net::{TcpListener, TcpStream};
use std::collections::HashMap;
use log::{info, warn};

// tcp协议的P2P节点实现
#[derive(Debug)]
struct TCPNode<'a> {
    cfg : &'a P2PConf,
    peers: HashMap<String,TcpStream>
}

impl<'a, 'b: 'a> HostStuffModule<'b> for TCPNode<'a> {
    type Service = Self;
    fn new(cfg :&'b HotStuffConf) -> Self {
        info!("创建TCPNode作为P2P服务节点");
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

#[cfg(test)]
mod tests{

    use log::info;
    use super::TCPNode;
    use crate::HotStuffConf;
    use crate::HostStuffModule;
    #[test]
    fn tcpnode_new_work() {
        log::set_max_level(log::LevelFilter::Debug);
        let cfg = HotStuffConf::new().unwrap();
        let node  = TCPNode::new(&cfg);
        info!("tcp node: {:?}", node);
        // println!("tcp node: {:?}", node);
    }

}