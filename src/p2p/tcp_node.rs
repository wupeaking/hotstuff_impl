use crate::HostStuffModule;
use crate::config::HotStuffConf;

// tcp协议的P2P节点实现
struct TCPNode {
}

impl HostStuffModule for TCPNode {
    type service = Self;
    fn new(cfg :&HotStuffConf) -> Self {
        return TCPNode{};
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
