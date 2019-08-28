#include "parser.hpp"

#include <iostream>

#include <boost/algorithm/string/split.hpp>
#include <boost/algorithm/string/erase.hpp>
#include <boost/algorithm/string/predicate.hpp>
#include <boost/program_options/variables_map.hpp>

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
        boost::algorithm::erase_all(code, "\n");
        boost::algorithm::erase_all(code, "\r");
        recurseNodes(code, parent);
        return parent;
    }

    void Parser::recurseNodes(std::string const& code, std::weak_ptr<AbstractNode> const& parent, int depth) {
        using nodeType = std::function<std::shared_ptr<AbstractNode>(std::string const&, AbstractNode::ParentRef)>;
        std::map<std::string, nodeType> nameToNode;
        nameToNode.emplace("package", [](auto name, auto parent) { return std::make_shared<PackageNode>(name, parent); });
        nameToNode.emplace("default", [](auto name, auto parent) { return std::make_shared<ParseNode>(name, parent); });
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
                    boost::algorithm::starts_with(blockWithInfo, "package");
                    std::cout << blockWithInfo << std::endl;
//                    auto child = std::make_shared<ParseNode>(std::move(blockWithInfo), parent);
                    auto child = nameToNode["default"](blockWithInfo, parent);
                    parent.lock()->addChild(child);
                    std::string blockContents = code.substr(blockStart + 1, i - blockStart - 1);
                    recurseNodes(blockContents, child, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
