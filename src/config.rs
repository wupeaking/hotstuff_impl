use crate::p2p::Peer;
pub struct HotStuffConf {
    pub p2p_cfg :P2PConf,
}

pub struct P2PConf {
    listen_addr :String,
    peers : Vec<Peer>,
} 