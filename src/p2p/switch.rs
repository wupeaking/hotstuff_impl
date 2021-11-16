use anyhow::Error;
use serde::{Deserialize, Serialize};

//type OnReceive func(modelID string, msgBytes []byte, p *Peer)

pub type OnReceive = fn(model_id: String, msg: &[u8], p: &Peer);

pub trait SwitcherI {
    fn broadcast(&self, model_id: String, msg: &BroadcastMsg) -> Error;
    fn broadcast_to_peer(&self, model_id: String, msg: &BroadcastMsg, p: &Peer) -> Error;
    fn broadcast_except_peer(&self, model_id: String, msg: &BroadcastMsg, p: &Peer) -> Error;
    fn remove_peer(&mut self, p: &Peer) -> Error;

    fn register_receive_callback(&mut self, model_id: &str, call_back: OnReceive) -> Error;
}

#[derive(Serialize, Deserialize)]
pub struct BroadcastMsg {
    pub model_id: String,
    pub msg: Vec<u8>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Peer {
    pub id: String,      // 定义peerid  每个peerid应该是唯一的
    pub address: String, // 地址
}