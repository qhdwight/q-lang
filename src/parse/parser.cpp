#include "parser.hpp"

#include <iostream>

#include <boost/algorithm/string/erase.hpp>
#include <boost/program_options/variables_map.hpp>

#include <util/read.hpp>
#include <parse/node/parse_node.hpp>

namespace ql::parse {
    std::shared_ptr<MasterNode> Parser::parse(po::variables_map& options) {
        auto sources = options["input"].as<std::vector<std::string>>();
        std::string sourceFileName = sources[0];
        auto src = util::readAllText(sourceFileName);
        auto node = getNodes(src.value());
        return node;
    }

//    std::vector<std::string>
//    std::vector<std::string>

    std::shared_ptr<MasterNode> Parser::getNodes(std::string code) {
        auto parent = std::make_shared<MasterNode>();
        boost::algorithm::erase_all(code, "\n");
        boost::algorithm::erase_all(code, "\r");
        recurseNodes(code, parent);
        return parent;
    }

    void Parser::recurseNodes(const std::string& code, const std::weak_ptr<AbstractNode>& parent, int depth) {
        auto level = 0;
        int blockInfoStart = 0, blockStart = 0;
        for (int i = 0; i < static_cast<int>(code.size()); i++) {
            char c = code[i];
            if (c == '{') {
                if (level++ == 0) {
                    blockStart = i;
                }
            } else if (c == '}') {
                if (--level == 0) {
                    std::string blockWithInfo = code.substr(blockInfoStart, i - blockInfoStart + 1);
                    std::cout << blockWithInfo << std::endl;
                    auto child = std::make_shared<ParseNode>(std::move(blockWithInfo), parent);
                    parent.lock()->addChild(child);
                    std::string blockContents = code.substr(blockStart + 1, i - blockStart - 1);
                    recurseNodes(blockContents, child, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
