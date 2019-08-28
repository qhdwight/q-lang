#include "parser.hpp"

#include <iostream>

#include <boost/algorithm/string/trim.hpp>
#include <boost/algorithm/string/split.hpp>
#include <boost/range/algorithm_ext/erase.hpp>
#include <boost/algorithm/string/predicate.hpp>
#include <boost/program_options/variables_map.hpp>
#include <boost/algorithm/string/classification.hpp>

#include <util/read.hpp>
#include <parser/node/package_node.hpp>

namespace ql::parser {
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
//        boost::erase_all(code, "\n");
//        boost::erase_all(code, "\r");
        boost::remove_erase_if(code, boost::is_any_of("\n\r"));
        recurseNodes(code, parent);
        return parent;
    }

    void Parser::recurseNodes(std::string const& code, std::weak_ptr<AbstractNode> const& parent, int depth) {
        using nodeFunc = std::function<std::shared_ptr<AbstractNode>(std::string&&, std::vector<std::string>&&, AbstractNode::ParentRef)>;
        std::map<std::string, nodeFunc> nameToNode;
        nameToNode.emplace("pckg", [](auto name, auto tokens, auto parent) {
            return std::make_shared<PackageNode>(std::move(name), std::move(tokens), parent);
        });
        nameToNode.emplace("default", [](auto name, auto tokens, auto parent) {
            return std::make_shared<ParseNode>(std::move(name), std::move(tokens), parent);
        });
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
                    auto delimiters = boost::is_any_of("\t ");
                    // Trim is necessary since split will include empty strings in beginning if we do not
                    boost::trim_if(blockWithInfo, delimiters);
                    // Split by tabs and spaces into tokens, which we use to find what type of node to create
                    std::vector<std::string> tokens;
                    boost::split(tokens, blockWithInfo, delimiters, boost::token_compress_on);
                    std::string const& nodeName = tokens[0];
                    std::cout << blockWithInfo << std::endl;
                    // Check if we have a generator function that can make this requested time, or else use default
                    auto it = nameToNode.find(nodeName);
                    nodeFunc blockNodeFunc = it == nameToNode.end() ? nameToNode["default"] : it->second;
                    auto child = blockNodeFunc(std::move(blockWithInfo), std::move(tokens), parent);
                    // Add children to parent node, parent node is owning via a shared pointer
                    parent.lock()->addChild(child);
                    // Recurse on the inner contents of the block so that each node added is for one block only
                    std::string blockContents = code.substr(blockStart + 1, i - blockStart - 1);
                    recurseNodes(blockContents, child, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
