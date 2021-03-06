use crate::p2p::Peer;
use anyhow::Error;
use serde::{Deserialize, Serialize};
use serde_yaml;

#[derive(Serialize, Deserialize)]
pub struct HotStuffConf {
    pub p2p_cfg: P2PConf,
    pub node_id: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct P2PConf {
    pub listen_addr: String,
    peers: Vec<Peer>,
}

impl HotStuffConf {
    pub fn new() -> Result<Box<Self>, Error> {
        let f = std::fs::File::open("src/config.yaml")?;
        let cfg = serde_yaml::from_reader(f)?;
        Ok(cfg)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cfg_from_yaml() {
        HotStuffConf::new().unwrap_or_else(|e| {
            panic!("{:?}", e);
        });
    }

    // #[test]
    fn cfg_to_yaml() {
        let scrape_config = HotStuffConf {
            p2p_cfg: P2PConf {
                listen_addr: "127.0.0.".to_owned(),
                peers: vec![Peer {
                    id: "11111".to_owned(),
                    address: "22.22".to_owned(),
                }],
            },
            node_id: "11111".to_string(),
        };

        let f = std::fs::OpenOptions::new()
            .write(true)
            .create(true)
            .open("new_config.yml")
            .expect("Couldn't open file");
        serde_yaml::to_writer(f, &scrape_config).unwrap();
    }
}
