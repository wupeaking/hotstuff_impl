use crate::config::{HotStuffConf, P2PConf};
use crate::HostStuffModule;
use anyhow::{Error, anyhow};
use log::{info, warn};
use std::collections::HashMap;
use tokio::net::{TcpListener, TcpStream};

// tcp协议的P2P节点实现
#[derive(Debug)]
struct TCPNode<'a> {
    cfg: &'a P2PConf,
    peers: HashMap<String, TcpStream>,
}

impl<'a, 'b: 'a> HostStuffModule<'b> for TCPNode<'a> {
    type Service = Self;
    fn new(cfg: &'b HotStuffConf) -> Self {
        info!("创建TCPNode作为P2P服务节点");
        return TCPNode {
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

impl<'a> TCPNode<'a>{
    async fn try_get_connected_peer(&mut self, peer_id:&str) -> Result<&mut TcpStream, anyhow::Error> {
        if let Some(v) = self.peers.get_mut(peer_id){
            return Ok(v);
        }

        let peer = self.peers.iter().find(|x|x.0==peer_id);
        let peer_addr = match peer {
            Some(v) =>v.0,
            _ => {
                return Err(anyhow!("未知的peer"));
            }
        };
        // Connect to a peer
        let mut stream = TcpStream::connect(peer_addr).await?;

        self.peers.insert(peer_id.to_string(), stream);

        // // Write some data.
        // stream.write_all(b"hello world!").await?;
        Ok(&mut stream)
    }
}

#[cfg(test)]
mod tests {

    use super::TCPNode;
    use crate::HostStuffModule;
    use crate::HotStuffConf;
    use log::info;
    #[test]
    fn tcpnode_new_work() {
        log::set_max_level(log::LevelFilter::Debug);
        let cfg = HotStuffConf::new().unwrap();
        let node = TCPNode::new(&cfg);
        info!("tcp node: {:?}", node);
        // println!("tcp node: {:?}", node);
    }
}
