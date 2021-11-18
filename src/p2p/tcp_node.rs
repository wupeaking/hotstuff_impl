use crate::config::{HotStuffConf, P2PConf};
use crate::HostStuffModule;
use anyhow::{anyhow, Error};
use log::{info, warn};
use std::collections::HashMap;
use tokio::net::{TcpListener, TcpStream};

use super::{BroadcastMsg, OnReceive, Peer, SwitcherI};

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
        let addr = self.cfg.listen_addr.clone();
        use std::sync::Arc;
        let this = Arc::new(self);
        tokio::spawn(async move {
            let listener = TcpListener::bind(addr).await.unwrap();
            loop {
                let (socket, _) = listener.accept().await.unwrap();
                this.process_stream(socket);
            }
        });

        Ok(())
    }

    // todo::
    fn stop(&mut self) -> Result<(), anyhow::Error> {
        Ok(())
    }
}

impl<'a> TCPNode<'a> {
    async fn try_get_connected_peer(
        &mut self,
        peer_id: &str,
    ) -> Result<&mut TcpStream, anyhow::Error> {
        // 首先获取可变peers 和mut self隔离开 不然编译器又报self多次引用了
        let mut_peers = &mut self.peers;
        let mut exit = false;
        if let Some(_) = mut_peers.get_mut(peer_id) {
            exit = true;
        }

        // 说明没有获取到
        if !exit {
            let peer = mut_peers.into_iter().find(|x| x.0 == peer_id);
            let peer_addr = match peer {
                Some(v) => v.0,
                _ => {
                    return Err(anyhow!("未知的peer"));
                }
            };
            // Connect to a peer
            let stream = TcpStream::connect(peer_addr).await?;

            mut_peers.insert(peer_id.to_string(), stream);
        }
        // // Write some data.
        // stream.write_all(b"hello world!").await?;
        let v = mut_peers.get_mut(peer_id).unwrap();

        Ok(v)
    }

    async fn process_stream(&self, stream: TcpStream) {

    }
}

impl<'a> SwitcherI for TCPNode<'a> {
    fn broadcast(&self, model_id: String, msg: &BroadcastMsg) -> Result<(), Error> {
        Ok(())
    }
    fn broadcast_to_peer(
        &self,
        model_id: String,
        msg: &BroadcastMsg,
        p: &Peer,
    ) -> Result<(), Error> {
        Ok(())
    }
    fn broadcast_except_peer(
        &self,
        model_id: String,
        msg: &BroadcastMsg,
        p: &Peer,
    ) -> Result<(), Error> {
        Ok(())
    }
    fn remove_peer(&mut self, p: &Peer) -> Result<(), Error> {
        Ok(())
    }

    fn register_receive_callback(
        &mut self,
        model_id: &str,
        call_back: OnReceive,
    ) -> Result<(), Error> {
        Ok(())
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

    #[test]
    fn tcpnode_start_work() {
        log::set_max_level(log::LevelFilter::Debug);
        let cfg_box = HotStuffConf::new().unwrap();

        let mut node = TCPNode::new(&cfg_box);

        node.start().unwrap();
    }
}
